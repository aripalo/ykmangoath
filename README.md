<h1><img src="/assets/ykmangoath.svg" height="32px" alt="logo" /> <code>ykmangoath</code> - Ykman OATH TOPT with Go</h1>

🚧  **Work-in-Progress**

Yet another **Go wrapper for [`ykman`](https://developers.yubico.com/yubikey-manager/) for generating [`OATH TOPT`](https://en.wikipedia.org/wiki/Time-based_one-time_password) codes from Yubikey**.

There are already [some](https://github.com/99designs/aws-vault/blob/master/prompt/ykman.go) [packages](https://github.com/joshdk/ykmango) out there which already wrap Yubikey CLI – `ykman` – for Go to manage OATH TOPT, but they lack all or some of the following features:

- Go Context support (handy for timeouts/cancellation etc)
- Multiple Yubikeys (identified by _device ID_)
- Password protected Yubikey OATH applications

Hence, this package, which covers those features! Big thanks to [`joshdk/ykmango`](https://github.com/joshdk/ykmango) and [`99designs/aws-vault`](https://github.com/99designs/aws-vault/blob/master/prompt/ykman.go) as they heavily influenced this library. Also this library is somewhat based on the previous implementation of Yubikey support in [`aripalo/vegas-credentials`](https://github.com/aripalo/vegas-credentials) (which this partly replaces in near future).

<br/>

## Installation

Requires:
- Yubikey Series 5 device (or newer with a `OATH TOPT` support)
- [`ykman`](https://developers.yubico.com/yubikey-manager/) Yubikey Manager CLI
- Go `1.18` or newer (for development)

```sh
go get github.com/aripalo/ykmangoath
```

<br/>

## Usage

### Initialization

This `ykmangoath` library provides a struct `OathAccounts` which represents a the main functionality of Yubikey OATH accounts (via `ykman` CLI). You can “create an instance” of the struct with `ykmangoath.New` and provide the following:
- Context (type of `context.Context`) which allows you to implement for example cancellations and timeouts
- Device Serial Number which is the 8+ digit serial number of your Yubikey device which you can find:
  - Printed in the back of your physical Yubikey device
  - By running command `ykman info` in your terminal

```go
package main

import (
	"github.com/aripalo/ykmangoath"
)

func main() {
	myTimeout := 20*time.Second
	ctx, cancel := context.WithTimeout(context.Background(), myTimeout)
	defer cancel()

	deviceSerial := "12345678" // can be empty string if you only use one Yubikey device
	oathAccounts := ykmangoath.New(ctx, deviceSerial)
}
```

Once initialized, you may perform operations on it such as [`List`](#list-accounts) or [`Code`](#request-code) methods. See [Managing Password](#managing-password) if your Yubikey OATH application is password protected.

<br/>

### List Accounts

Implements `ykman --device 12345678 oath accounts list` with Go.

```go
accounts, err := oathAccounts.List()
if err != nil {
  log.Fatal(err)
}

fmt.Println(accounts)
```

<br/>

### Request Code

Implements `ykman --device 12345678 oath accounts code --single '<issuer>:<name>'` with Go.

```go
account := "<issuer>:<name>"
code, err := oathAccounts.Code(account)
if err != nil {
  log.Fatal(err)
}

fmt.Println(code)
```

<br/>

### Managing Password

#### Direct Assign

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

#### Prompt Function

Instead of assigning the password directly ahead-of-time, you may provide a **_password prompt function_ that will be executed only if password is required**.

It must return a password `string` (which can be empty) and an `error` (which of course could be `nil` on success). The password prompt function will also receive the `context.Context` given in `ykmangoath.New` initialization, therefore your password prompt function can be cancelled (for example due to timeout).

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	return "p4ssword", nil
}

err := oathAccounts.SetPasswordPrompt(myPasswordPrompt)
```

##### Retrieve the prompted password

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	return "p4ssword", nil
}

err := oathAccounts.SetPasswordPrompt(myPasswordPrompt)
// handle err

code, err := oathAccounts.Code("<issuer>:<name>")
// handle err
// do something with code

password, err := oathAccounts.GetPassword()
// handle err
// do something with password (e.g. cache it somewhere)
```


This can be useful if you wish to cache the Yubikey OATH application password for short periods of time in your own application. How you cache it (hopefully somewhat securely) is up to you.
