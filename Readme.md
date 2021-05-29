# Resource MD5
This is a tool that prints MD5 hash for each given resource.

It can fetch resources in parallel. By default, it creates 10 workers. You can control this behavior by passing `-parallel` flag.

Program will return an error and halt in the following cases:
- HTTP request has been failed
- no scheme provided for an URL
- server returned non-200 status code
- flags parsing failed

## Usage
```shell
go build -o rmd5 ./main.go
./rmd5 -parallel 5 https://url1.com http://url2.com

```

Note: every URL should have a scheme, e.g. https:// or http://.
