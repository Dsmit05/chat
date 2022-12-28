# chat
example chat rooms on websocket

### endpoint

`localhost:8080/` - info

`localhost:8080/rooms/{id}` - websocket chat rooms

### msg format
```json
{
  "text": ""
}
```

### run
1. `git clone https://github.com/Dsmit05/chat.git`
2. `go run cmd\chat\main.go`

or

1. `git clone https://github.com/Dsmit05/chat.git`
2. `docker-compose up`
