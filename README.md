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

	oath := ykmangoath.New(ctx, "12345678")

	accounts, err := oath.List()
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

	oath := ykmangoath.New(ctx, "12345678")

	code, err := oath.Code("<issuer>:<name>")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}
```

<br/>

### Managing Password

#### Direct configuration

```go
oath := ykmangoath.New(ctx, "12345678")
err := oath.SetPassword("p4ssword")
```

#### Prompt Function

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	return "p4ssword", nil
}

oath := ykmangoath.New(ctx, "12345678")
err := oath.SetPasswordPrompt(myPasswordPrompt)
```

##### Retrieve the prompted password

```go
func myPasswordPrompt(ctx context.Context) (string, error) {
	return "p4ssword", nil
}

oath := ykmangoath.New(ctx, "12345678")
err := oath.SetPasswordPrompt(myPasswordPrompt)
// handle err

code, err := oath.Code("<issuer>:<name>")
// handle err
// do something with code

password, err := oath.GetPassword()
// handle err
// do something with password (e.g. cache it somewhere)
```


This can be useful if you wish to cache the Yubikey OATH application password for short periods of time in your own application. How you cache it (hopefully somewhat securely) is up to you.
