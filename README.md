# Test task

Design and implement “Word of Wisdom” tcp server.

- TCP server should be protected from DDOS attacks with the [Prof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, the server should send one of the quotes from “Word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge


## Run

```bash
cd server
docker build --tag hashcash-server .
docker run -p 8000:8000 -d hashcash-server
cd ../client
docker build --tag hashcash-client .
docker run -d hashcash-client
```

The server uses the [Hashcash](https://en.wikipedia.org/wiki/Hashcash) proof of work algorithm with complexity equals to 3.

## Client output
```bash
Got challenge 172.17.0.1 1650825229762463703
Got solution 00006996f1fc573c8587a7b71e70ed2cb43e2444 after 144032 steps
“The best way out is always through“ - Robert Frost
“Always Do What You Are Afraid To Do” – Ralph Waldo Emerson
“The best way out is always through“ - Robert Frost
“The journey of a thousand miles begins with one step.” – Lao Tzu
“The best way out is always through“ - Robert Frost
```

## Why Hashcash?

It uses a commonly known `sha1` hash function, because it will not take much computational resources for server to check sha1, thus adding DDOS hashcash protection will not overload it much. While on client side it will take significant computational resources to connect to server, thus protect server from DDOS attacks.  
