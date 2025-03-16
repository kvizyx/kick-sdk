# About the Kick SDK

Kick SDK is a toolkit for Golang applications that provides a clean way to interact with Kick public APIs.

## Installation

```bash
go get github.com/glichtv/kick-sdk
```

## Documentation
This documentation only belongs to the Kick SDK, if you are looking for the official documentation you can find it
[here](https://docs.kick.com/).

## Unofficial Endpoints

It is possible that support for closed (unofficial) endpoints will be added in the future, but for now only publicly
open and documented APIs are supported, as one of the goals of the Kick SDK is stability, which is hard to achieve with
additional support for unofficial endpoints.

However, Kick SDK allows you to change client's default base URLs and do custom requests, so in fact you can effortlessly
implement any endpoint by yourself if you need to.
