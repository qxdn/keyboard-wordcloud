# keyboard-wordcloud

only support windows now

record your keyboard input and generate wordclouds everyday

the default output at `"~/Pictures/wordclouds"`

# build
run follwing code
```
go mod download
// with comman line
go build -tag windows
// without command line
go build -tag windows windows -ldflags -H=windowsgui
```

# TODO

- [ ] use config.yaml or config.json
- [ ] windows 托盘