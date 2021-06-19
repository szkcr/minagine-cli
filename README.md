# automate-minagine

A CLI tool for dakoku MINAGINE

## build
```
$ go build -o automate-minagine
```

## usage
```
$ automate-minagine -h
Usage:
  automate-minagine [OPTIONS]

Application Options:
  -d, --domain=                   domain of your account (tenant)
  -u, --user=                     user id of your account
  -p, --password=                 password of your account
  -a, --action=[checkin|checkout] desired action
  -h, --webhook=                  [option] webhook url for reporting action result

Help Options:
  -h, --help                      Show this help message
```
