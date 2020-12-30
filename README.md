GET Today Todo List FROM Asana

```go
go get -v
go run main.go
```

```sh
touch .env
```

.env 파일 
```
PORT=1324
ACCESS_TOKEN=
```

url에 uid 설정 개인마다 바꿔서 가야함. ACCESS_TOKEN은 개인마다 발급해야하는 듯.


## TODO

- 유저마다 ACCESS_TOKEN, UID 환경변수 등록
- SLACK, SLASH COMMAND 연동
