# csv2ss
csv to google spreadsheets

## Install

```
$ go get github.com/itsubaki/csv2ss
$ cd ${GOPATH}/src/github.com/itsubaki/csv2ss
$ make
```

## Example

```
$ export GOOGLE_CREDENTIALS_PATH=test/secret/credentials.json
$ export GOOGLE_TOKEN_PATH=test/secret/token.json
$ cat test/test.csv | csv2ss
```
