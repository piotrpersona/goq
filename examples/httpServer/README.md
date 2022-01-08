# Http Server save every request example

## Run

Start server:
```console
go run main.go
```

Server output:

```console
Server listening at :8080
```

Send request:
```console
curl localhost:8080
```
Console output:
```console
request saved in 3.494µs⏎
```

Server output:
```console
2022-01-08T11:54:16Z {Method:GET URL:/ Proto:HTTP/1.1 ProtoMajor:1 ...
```
