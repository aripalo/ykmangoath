<div align="center">
	<br/>
	<br/>
	<img width="200" src="assets/ykmangoath.svg" alt="Got" />
  <h1>
  <code>ykmangoath</code>
  <br/>
  <span>Ykman OATH TOPT with Go</span>
  </h1>
  <div>

  _Yet another_ [`ykman`](https://developers.yubico.com/yubikey-manager/) Go lib for requesting [OATH TOPT](https://en.wikipedia.org/wiki/Time-based_one-time_password) Multi-Factor Authentication Codes from [Yubikey](https://www.yubico.com/products/yubikey-5-overview/) Devices.

  </div>
  <hr/>
  <br/>
</div>


🚧  **Work-in-Progress**

There are already [some](https://github.com/99designs/aws-vault/blob/master/prompt/ykman.go) [packages](https://github.com/joshdk/ykmango) out there which already wrap Yubikey CLI for Go to manage OATH TOPT, but they lack all or some of the following features:

- **Go Context support** (handy for timeouts/cancellation etc)
- **Multiple Yubikeys** (identified by Device Serial Number)
- **Password protected** Yubikey OATH applications

Hence, this package, which covers those features! Big thanks to [`joshdk/ykmango`](https://github.com/joshdk/ykmango) and [`99designs/aws-vault`](https://github.com/99designs/aws-vault/blob/master/prompt/ykman.go) as they heavily influenced this library. Also this library is somewhat based on the previous implementation of Yubikey support in [`aripalo/vegas-credentials`](https://github.com/aripalo/vegas-credentials) (which this partly replaces in near future).

This library supports only a small subset of features of Yubikeys & OATH account management, this is [by design](#design).

<br/>

## Installation

### Requirements

- Yubikey Series 5 device (or newer with a `OATH TOPT` support)
- [`ykman`](https://developers.yubico.com/yubikey-manager/) Yubikey Manager CLI as a runtime dependency
- Go `1.18` or newer (for development)

### Get it

```sh
go get -u github.com/aripalo/ykmangoath
```

<br/>

## Getting Started

This `ykmangoath` library provides a struct `OathAccounts` which represents a the main functionality of Yubikey OATH accounts (via `ykman` CLI). You can “create an instance” of the struct with `ykmangoath.New` and provide the following:
- Context (type of `context.Context`) which allows you to implement things like cancellations and timeouts
- Device Serial Number which is the 8+ digit serial number of your Yubikey device which you can find:
  - printed in the back of your physical Yubikey device
  - by running command `ykman info` in your terminal

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aripalo/ykmangoath"
)

func main() {
	myTimeout := 20*time.Second
	ctx, cancel := context.WithTimeout(context.Background(), myTimeout)
	defer cancel()

	deviceSerial := "12345678" // can be empty string if you only use one Yubikey device
	oathAccounts := ykmangoath.New(ctx, deviceSerial)

	// after this you can perform various operations on oathAccounts...
}
```

Once initialized, you may perform operations on `oathAccounts` such as [`List`](#list-accounts) or [`Code`](#request-code) methods. See [Managing Password](#managing-password) if you wish to support password protected Yubikey OATH applications (you really should).

<br/>

## List Accounts

List OATH accounts configured in the Yubikey device:

```go
accounts, err := oathAccounts.List()
if err != nil {
  log.Fatal(err)
}

fmt.Println(accounts)
```

The above is the same as running the following in your terminal:
```sh
ykman --device 12345678 oath accounts list
```

Example Go output:
```go
[
  "Amazon Web Services:john.doe@example",
]
```

<br/>

## Request Code

Requests a _Time-based One-Time Password_ (TOPT) 6-digit code for given account (such as "<issuer>:<name>") from Yubikey OATH application.

```go
account := "<issuer>:<name>"
code, err := oathAccounts.Code(account)
if err != nil {
  log.Fatal(err)
}

fmt.Println(code)
```

The above is the same as running the following in your terminal:
```sh
ykman --device 12345678 oath accounts code --single '<issuer>:<name>'
```

Example Go output:
```go
"123456"
```

<br/>

## Managing Password

An end-user with Yubikey device may wish to password protect the Yubikey's OATH application. Generally speaking this is a good idea as it adds some protection from theft: A bad actor with someone else's Yubikey device can't actually use the device to generate TOPT MFA codes unless they somehow also know the device password.

The password protection for the Yubikey device's OATH application can be set either via the Yubico Authenticator GUI application or via command-line with `ykman` by running:
```sh
ykman oath access change
```

But, if the device is configured with a password protected OATH application, it means that there needs to be a way to provide the password for `ykmangoath`: Luckily this is one of the benefits of this specific library as it supports just that by either:
- [Directly Assigning](#direct-assign) the password ahead-of-time
- [Prompt Function](#prompt-function) which you can use to ask the password from end-user

There's also functionality to [retrieve the prompted password](#retrieve-the-prompted-password) given by end-user, so you may implement caching mechanisms to provide a smoother user experience where the end-user doesn't have to type in the password for every Yubikey operation; Remember there are of course tradeoffs with security vs. user experience with caching the password!

## Direct Assign

```go
err := oathAccounts.SetPassword("p4ssword")
// handle err
account := "<issuer>:<name>"
code, err := oathAccounts.Code(account)
// handle err
```

The above is the same as running the following in your terminal:
```sh
ykman --device 12345678 oath accounts code --single '<issuer>:<name>' --password 'p4ssword'
```

## Prompt Function

Instead of assigning the password directly ahead-of-time, you may provide a **_password prompt function_ that will be executed only if password is required**. Often you'll use this to actually ask the password from end-user – either via terminal `stdin` or by showing a GUI dialog with tools such as [`ncruces/zenity`](https://github.com/ncruces/zenity).

It must return a password `string` (which can be empty) and an `error` (which of course could be `nil` on success). The password prompt function will also receive the `context.Context` given in `ykmangoath.New` initialization, therefore your password prompt function can be cancelled (for example due to timeout).

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	password := "p4ssword" // in real life, you'd resolve this value by asking the end-user
	return password, nil
}

err := oathAccounts.SetPasswordPrompt(myPasswordPrompt)
```

### Retrieve the prompted Password

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	password := "p4ssword" // in real life, you'd resolve this value by asking the end-user
	return password, nil
}

err := oathAccounts.SetPasswordPrompt(myPasswordPrompt)
// handle err

code, err := oathAccounts.Code("<issuer>:<name>")
// handle err
// do something with code

password, err := oathAccounts.GetPassword()
// handle err
// do something with password (e.g. cache it somewhere):
myCacheSolution.Set(password) // ... just an example
```


This can be useful if you wish to cache the Yubikey OATH application password for short periods of time in your own application, so that the user doesn't have to type in the password everytime (remember: the physical touch of the Yubikey device should be the _actual_ second factor). How you cache it (hopefully somewhat securely) is up to you.

<br/>

## Design

### This tool is designed only for _retrieval_ of specific information from a Yubikey device:

- List of configured OATH accounts
- _Time-based One-Time Password_ (TOPT) code for given OATH account
- ... with a support for password protected Yubikey OATH applications

### By design, this tool does NOT support:

- Setting or changing the Yubikey OATH application password: Configuring the password should be an explicit end-user operation (which they can do via Yubico Authenticator GUI or `ykman` CLI) – We do not want to enable situations where this library renders end-user's Yubikey OATH accounts useless by setting a password unknown to the end-user.
- Removing or renaming OATH accounts from the Yubikey device: This is another area where – either by accident or on purpose – one could harm the end-user.

If you need some of the above-mentioned unsupported features, feel free to implement them yourself – for example by forking this library but we will never accept a pull request back into this project which implements those features.

## This tool _may_ implement following features in the future:

- https://github.com/aripalo/ykmangoath/issues/1 Adding new OATH acccounts
- https://github.com/aripalo/ykmangoath/issues/2 Password protecting a Yubikey device's OATH application given that:
    1. there are no OATH accounts configured (i.e. the device OATH application is empty/unused)
    2. the device OATH application does not yet have a password protection enabled

But this is not a promise to implement them. If you feel like it, you can create a Pull Request!
