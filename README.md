# gomparator 

gomparator is used for comparing HTTP JSON responses of two hosts by checking if they respond with the same json (deep equal ignoring order) and status code.

Note: Only supports HTTP GET methods. More verbs are yet to come.

## Download and install

    go get -u github.com/emacampolo/gomparator

## Create a file with relative URL 

    eg:
    
    /v1/payment_methods?client.id=1
    /v1/payment_methods?client.id=2
    /v1/payment_methods?client.id=3
    
## Run

```sh
$ gomparator -path "/path/to/file/with/urls" -host "http://host1.com" -host "http://host2.com" -header "X-Auth-Token: abc"
```
![](example.gif)

By Default it will use 1 worker and the rate limit will be 5 req/s. This can be overridden. For more info run:

```sh
$ gomparator -h
```
