# minagine-cli

A CLI tool for dakoku MINAGINE

## build

```
$ go build -o minagine-cli
```

## usage

```
$ minagine-cli -h
Usage:
  minagine-cli [OPTIONS]

Application Options:
  -d, --domain=                   domain of your account (tenant)
  -u, --user=                     user id of your account
  -p, --password=                 password of your account
  -a, --action=[checkin|checkout] action to be performed
  -w, --webhook=                  [option] webhook url for reporting action result

Help Options:
  -h, --help                      Show this help message
```
