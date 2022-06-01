# Draw - multiplayer sketchpad

Inspired by @mrdoob :)

I wanted to see how much better performance I could get rewriting the server in Go.

## Features

- Live user count
- Unlimited colours

## How to run

```
go mod download
go run .
```

or with Docker

```
docker build -t draw .
docker run -d -p 3000:3000 --name draw draw

//to stop
docker stop draw
```

