json2env
=======

[![Test Status](https://github.com/dekokun/json2env/workflows/test/badge.svg?branch=master)][actions]
[![Coverage Status](https://coveralls.io/repos/dekokun/json2env/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/dekokun/json2env?status.svg)][godoc]

[actions]: https://github.com/dekokun/json2env/actions?workflow=test
[coveralls]: https://coveralls.io/r/dekokun/json2env?branch=master
[license]: https://github.com/dekokun/json2env/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/dekokun/json2env

json2env execute commands with environment variables made from JSON.

## Synopsis

```go
$ echo '{"key":"value"}' | json2env --keys "key" /path/to/command [...]
```

```go
$ echo '{"examplekey1":"value1", "examplekey2":"value2", "examplekey3":"value3"}' | json2env --keys "examplekey1,examplekey2" env | grep examplekey
examplekey1=value1
examplekey2=value2
```

## Motivation

I wanted to pass secret information as an environment variable to the ECS container in Fargate via the Secrets Manager, but I was unable to retrieve some data by specifying a key in the JSON of the Secrets Manager as shown below. Of course, you can get the entire JSON as an environment variable.

- https://docs.aws.amazon.com/AmazonECS/latest/userguide/specifying-sensitive-data-secrets.html#secrets-considerations
  - `For tasks that use the Fargate launch type, the following should be considered:`
  - `It is only supported to inject the full contents of a secret as an environment variable. Specifying a specific JSON key or version is not supported at this time.`

Since this feature is implemented in ECS on EC2, this feature will probably be implemented in Fargate in the future, so I wanted something to convert JSON to environment variables that can be used in the meantime.

And I made it mandatory to specify the name of the environment variable to use so that I can safely remove this tool when I no longer use it.

## Installation

```console
$ go get github.com/dekokun/json2env/cmd/json2env
```

## Author

[dekokun](https://github.com/dekokun)
