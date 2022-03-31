# `ykmangoath` - Ykman OATH TOPT with Go

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
- Go `1.18` or newer
- [`ykman`](https://developers.yubico.com/yubikey-manager/) Yubikey Manager CLI

```sh
go get github.com/aripalo/ykmangoath
```

<br/>

## Usage

### List Accounts

Implements `ykman --device 12345678 oath accounts list` with Go.

```go
package main

import (
	"github.com/aripalo/ykmangoath"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err, accounts := ykmangoath.List(ctx, ykmangoath.Options{DeviceID: "12345678"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(accounts)
}
```

<br/>

### Generate Code

Implements `ykman --device 12345678 oath accounts code --single '<issuer>:<name>'` with Go.

```go
package main

import (
	"github.com/aripalo/ykmangoath"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err, code := ykmangoath.Code(ctx, "<issuer>:<name>", ykmangoath.Options{DeviceID: "12345678"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}
```

<br/>

### Generate Code with Password Prompt

Implements `echo "p4ssword" | ykman --device 12345678 oath accounts code --single '<issuer>:<name>'` with Go.

If you have password protection for your Yubikey OATH application, you may either provide the password directly via `ykmangoath.Options` or provide your own _password prompt function_:

```go
package main

import (
	"github.com/aripalo/ykmangoath"
)

func myPasswordPrompt(ctx context.Context) (error, string) {
	return nil, "p4ssword"
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err, code := ykmangoath.CodeWithPasswordPrompt(
		ctx,
		myPasswordPrompt,
		"<issuer>:<name>",
		ykmangoath.Options{DeviceID: "12345678"},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}
```

This will first try the perform the operation without a password and if it detects a password is required, it will run the _password prompt function_ you provided (`myPasswordPrompt`) and try again with its result.

There's also a `ListWithPasswordPrompt` method to achieve the same password prompting functionality for [`List`](#list-accounts).

<br/>

### Retrieving the Password

There are also `ListWithPasswordPromptAndCache` and `CodeWithPasswordPromptAndCache` methods that contain a third return value: The password that succesfully unlocked the Yubikey OATH application:
```go
err, code, password := ykmangoath.CodeWithPasswordPromptAndCache(
	ctx,
	myPasswordPrompt,
	"<issuer>:<name>",
	ykmangoath.Options{DeviceID: "12345678"},
)
```

This can be useful if you wish to cache the Yubikey OATH application password for short periods of time in your own application. How you cache it (hopefully somewhat securely) is up to you.
