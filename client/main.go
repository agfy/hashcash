package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const numberOfTries = 1000000
const zerosNeeded = 3

func main() {
	var lastTimeUpdated time.Time
	var solution string

	for i := 0; i < 5; i++ {
		if !SolutionValid(lastTimeUpdated, solution) {
			challenge, err := GetChallenge()
			if err != nil {
				fmt.Println("failed to get challenge", err)
				return
			}
			fmt.Println("Got challenge", challenge)

			hash, err := ComputeSolution(challenge)
			if err != nil {
				fmt.Println("failed to find solution", err)
				continue
			}
			solution = challenge + " " + hash
			lastTimeUpdated = time.Now()
		}
		wow, err := GetWow(solution)
		if err != nil {
			fmt.Println("failed to get wow", err)
			continue
		}
		fmt.Println(wow)
	}
}

func SolutionValid(lastTimeUpdated time.Time, solution string) bool {
	if solution == "" || time.Now().Sub(lastTimeUpdated).Minutes() > 60 {
		return false
	}
	return true
}

func GetChallenge() (string, error) {
	resp, err := http.Get("http://172.17.0.1:8000/challenge")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	challenge, ok := data["Challenge"]
	if !ok {
		return "", errors.New("challenge missing")
	}
	return challenge.(string), nil
}

func ComputeSolution(challenge string) (string, error) {
	for i := 0; i < numberOfTries; i++ {
		solution := base64EncodeInt(i)
		hash := sha1Hash(challenge + " " + solution)
		if CheckHashValid(hash) {
			fmt.Println("Got solution", string(hash), "after", i, "steps")
			return solution, nil
		}
	}
	return "", errors.New("failed to find hash")
}

func base64EncodeInt(n int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(n)))
}

func CheckHashValid(hash []byte) bool {
	zeros := 0
	for _, symbol := range hash {
		if symbol != 48 { //zero symbol
			break
		}
		zeros++
	}
	if zeros <= zerosNeeded {
		return false
	}
	return true
}

func sha1Hash(s string) []byte {
	hash := sha1.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		return nil
	}
	strHash := fmt.Sprintf("%x", hash.Sum(nil))
	return []byte(strHash)
}

func GetWow(solution string) (string, error) {
	req, err := http.NewRequest("GET", "http://172.17.0.1:8000/wow", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Solution", solution)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	wow, ok := data["Wow"]
	if !ok {
		return "", errors.New("wow missing")
	}
	return wow.(string), nil
}
