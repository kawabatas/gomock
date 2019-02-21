# gomock

REST API mock server which returns dummy data.

## Setup
```
cp .envrc-example .envrc

go get github.com/rakyll/statik
statik -src=$PWD/mock_data

go build
```

## Server start
```
go build
./gomock -p [port]
```

curl example
```
curl http://localhost:8888/hello
curl -H 'X-FOO:BAR' http://localhost:8888/hello
curl http://localhost:8888/hello?foo=bar
curl -X POST http://localhost:8888/hello -d '{"foo":"bar"}'
```

## How to use in test
You can use this in `httptest` package.

See [example](https://github.com/kawabatas/gomock/tree/master/example).

## How to add dummy data
Add dummy data response (json) to `mock_data/` and add a call to `mock_data/routes.json`.

[rakyll/statik](https://github.com/rakyll/statik) makes it possible to distribute package with json file.

Please convert json file to Go binary. (Please update statik/statik.go.)
```
statik -src=$PWD/mock_data
```
