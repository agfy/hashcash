package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var wowLines []string

const nanosecondsInMinute = 1000 * 1000 * 1000 * 60
const zerosNeeded = 3

func main() {
	var err error
	rand.Seed(time.Now().UnixNano())
	wowLines, err = readFile("wow.txt")
	if err != nil {
		fmt.Println("failed to read file", err)
		return
	}

	r := gin.Default()
	r.GET("/wow", GetWowLine)
	r.GET("/challenge", GetChallenge)
	err = r.Run("0.0.0.0:8000")
	if err != nil {
		fmt.Println("error on main goroutine", err)
	}
}

func readFile(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func GetWowLine(c *gin.Context) {
	solution := c.GetHeader("Solution")
	if solution == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "reason": "empty Solution"})
		return
	}

	ip := c.ClientIP()
	err := Verify(solution, ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "reason": err.Error()})
		return
	}

	wowIndex := rand.Intn(len(wowLines))
	c.JSON(200, gin.H{"Wow": wowLines[wowIndex]})
}

func GetChallenge(c *gin.Context) {
	ip := c.ClientIP()
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	c.JSON(200, gin.H{"Challenge": ip + " " + now})
}

func Verify(response, ip string) error {
	vals := strings.Split(response, " ")
	if len(vals) != 3 {
		return errors.New("failed to parse response")
	}

	if vals[0] != ip {
		return errors.New("wrong ip")
	}

	unixNano, err := strconv.ParseInt(vals[1], 10, 64)
	if err != nil {
		return errors.New("failed to parse time")
	}
	minutesPassed := (time.Now().UnixNano() - unixNano) / nanosecondsInMinute
	if minutesPassed > 60 {
		return errors.New("response too old")
	}

	hash := sha1Hash(response)
	if hash == nil {
		return errors.New("failed to calculate hash")
	}
	if !CheckHashValid(hash) {
		return errors.New("not enough zeros")
	}
	return nil
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
