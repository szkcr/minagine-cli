# minagine-cli

A CLI tool for check-in/check-out MINAGINE

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
  -f, --force                     [option] skip pre-check of working day
  -w, --webhook=                  [option] webhook url for reporting action result

Help Options:
  -h, --help                      Show this help message
```
