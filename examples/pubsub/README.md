# Publish-Subscribe example

## Run

```console
go run main.go
```

Output:

```console
Start main goroutine #go: 1
goq.New will spawn another goroutine #go: 2
2 subscribers were spawned #go: 4
[topic-words]
[A B] <nil>
A: [{Key:1 Value:hello!} {Key:3 Value:abc!}]
B: [{Key:3 Value:abc!} {Key:1 Value:hello!}]
stop A
Group A unsubscribed #go: 3
A: [{Key:1 Value:hello!} {Key:3 Value:abc!}]
B: [{Key:3 Value:abc!} {Key:1 Value:hello!} {Key:4 Value:who's listening?}]
Close queue #go: 3
stop B
```
