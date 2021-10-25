# `teamhanko/hanko-go` Hanko API Client

This package is maintained by [Hanko](https://hanko.io).

## Contents
1. [Introduction](#introduction)
1. [Documentation](#documentation)
1. [Installation](#installation)
1. [Usage](#usage)
    1. [WebAuthn usage](#webauthn-usage)
        1. [Create a new Hanko API WebAuthn Client](#create-a-new-hanko-api-webauthn-client)
        1. [Register a WebAuthn credential](#register-a-webauthn-credential)
        1. [Authenticate with a registered WebAuthn credential](#authenticate-with-a-registered-webauthn-credential)
        1. [Making Transactions](#making-transactions)
        1. [Credential Management](#credential-management)
    1. [Passlink usage](#passlink-usage)
        1. [Create a new Hanko API Passlink Client](#create-a-new-hanko-api-passlink-client)
        1. [Passlink initialization](#passlink-initialization)
        1. [Passlink confirmation](#passlink-confirmation)
        1. [Passlink finalization](#passlink-finalization)
1. [Examples](#examples)
    1. [WebAuthn examples](#webauthn-examples)
        1. [Example of how to register credentials](#example-of-how-to-register-credentials)
        1. [Example of how to handle the authentication](#example-of-how-to-handle-the-authentication)
    1. [Passlink examples](#passlink-examples)
1. [Support](#support)

## Introduction
This repository contains an API client written in [Go](https://golang.org) that lets you communicate with the
[Hanko Authentication API](https://docs.hanko.io/overview)
to enable building passwordless authentication based on
[FIDO®](https://fidoalliance.org)/[WebAuthn](https://www.w3.org/TR/webauthn) or Passlinks (a.k.a. "login via email").

## Documentation

- [Hanko Docs](https://docs.hanko.io) website
- Hanko API Client [code documentation](https://pkg.go.dev/github.com/teamhanko/hanko-go)
- Hanko Authentication [API reference](https://docs.hanko.io/api/webauthn)

## Installation
1. Make sure [Go](https://golang.org) is installed.
2. Install the Hanko API Client:
```shell
$ go get -u github.com/teamhanko/hanko-go/webauthn
```
3. Import the client in your code:
```go
import "github.com/teamhanko/hanko-go/webauthn"
```

## Usage

### WebAuthn usage

#### Create a new Hanko API WebAuthn Client

First you need to log into the [Hanko Console](https://console.hanko.io). The console allows you to create a new
relying party, which also spins up a new Hanko Authentication instance (it's free, you can just try it 
out). To connect the client you need the
API URL of your instance, and it is also necessary to create a new API key in the settings of your relying party to
authenticate the client. For further information, see
[Getting started](https://docs.hanko.io/gettingstarted).

Once your instance is operational, and you obtained the API key (Key-ID and a secret) from the console, provide the
Key-ID via the [`WithHmac`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#Client.WithHmac) option 
and pass the secret directly into the constructor:

```go
var apiUrl string // e.g. "https://e7cd3792-dd57-4103-a6e3-4ac66fcd00f1.authentication.hanko.io"
var secret string // e.g. "FxCcAWGTEnT2P6mFbnQW4U..."

hankoWebAuthn = webauthn.NewClient(apiUrl, secret).
    WithHmac(hmacApiKeyId). // e.g. "66dc69ec-67dc-4901-aab3-a2bf4927671e"
    WithHttpClient(httpClient). // use a customized http client
    WithLogger(logger) // use a customized logger
```

#### Register a WebAuthn credential

Please visit [Hanko Docs](https://docs.hanko.io) to learn how a registration ceremony works and also
see the [example section](#examples) for a rough overview of how to utilize the client.

Registration initialization:
```go
var userId string // e.g. "65a3eba6-22cb-4c35-9881-b21fac6acfd0"
var userName string  // e.g. "alice@example.com"

user = webauthn.NewRegistrationInitializationUser(userId, userName).
    WithDisplayName(displayName) // e.g. "Alice"

authenticatorSelection = webauthn.NewAuthenticatorSelection().
    WithUserVerification(userVerification). // e.g. "required", "preferred", "discouraged"
    WithAuthenticatorAttachment(authenticatorAttachment). // e.g. "platform", "cross-platform"
    WithRequireResidentKey(requireResidentKey) // true or false

request = webauthn.NewRegistrationInitializationRequest(user).
    WithAuthenticatorSelection(authenticatorSelection).
    WithConveyancePreference(conveyancePreference) // e.g. "none", "indirect", "direct"

response, err = hankoWebAuthn.InitializeRegistration(request)
```
Registration finalization:
```go
// InitializeRegistration returns a RegistrationInitializationResponse that represents  
// PublicKeyCredentialCreationOptions that must be provided to the WebAuthn API to retrieve a 
// PublicKeyCredential, represented by a RegistrationFinalizationRequest.

request, err = webauthn.ParseRegistrationFinalizationRequest(registrationFinalizationRequest)
response, err = hankoWebAuthn.FinalizeRegistration(request)
```

#### Authenticate with a registered WebAuthn credential

Please visit [Hanko Docs](https://docs.hanko.io)  to learn how a authentication ceremony works and also
see the [example section](#examples) for a rough overview of how to utilize the 
client.

Authentication initialization:
```go
var userId string // e.g. "65a3eba6-22cb-4c35-9881-b21fac6acfd0"

user = webauthn.NewAuthenticationInitializationUser(userId)

request = webauthn.NewAuthenticationInitializationRequest().
    WithUser(user).
    WithUserVerification(userVerification).
    WithAuthenticatorAttachment(authenticatorAttachment)

response, err = hankoWebAuthn.InitializeAuthentication(request)
```
Authentication finalization:
```go
// InitializeAuthentication returns an AuthenticationInitializationResponse that represents  
// PublicKeyCredentialRequestOptions that must be provided to the WebAuthn API to retrieve a 
// PublicKeyCredential, represented by an AuthenticationFinalizationRequest.

request, err = webauthn.ParseAuthenticationFinalizationRequest(authenticationFinalizationRequest)
response, err = hankoWebAuthn.FinalizeAuthentication(request)
```

#### Making Transactions

A transaction is technically the equivalent of an authentication, with the difference that when initializing 
a transaction, a `transactionText` can be included, which is also used for signing the challenge.

Please visit [Hanko Docs](https://docs.hanko.io) for further information.

Transaction initialization:
```go
request = webauthn.NewTransactionInitializationRequest().
	WithTransaction(transactionText) // e.g. "Order #3242"
response, err = hankoWebAuthn.InitializeTransaction(request)
```
Transaction finalization:
```go
// InitializeTransaction returns a TransactionInitializationResponse that represents  
// PublicKeyCredentialRequestOptions that must be provided to the WebAuthn API to retrieve a 
// PublicKeyCredential, represented by a TransactionFinalizationRequest.

request, err = webauthn.ParseTransactionFinalizationRequest(transactionFinalizationRequest)
response, err = hankoWebAuthn.FinalizeTransaction(request)
```

#### Credential Management

Furthermore, the client offers the possibility to manage the registered credentials. If you create a productive 
application, make sure that you do not make these functions publicly available. Always verify whether the actor is
allowed to modify the credential.

```go
var credentialId string // e.g. "AQohBypyLBrx8R_UO0cWQuu7hhRGv7bPRRGtbQLrjl..."

// Get all details of the specified credential.
credential, err = hankoWebAuthn.GetCredential(credentialId)

// Update the name of a credential.
updateRequest = webauthn.NewCredentialUpdateRequest().
	WithName(newName) // e.g. "My Security Key"
credential, err = hankoWebAuthn.UpdateCredential(credentialId, updateRequest)

// Delete the specified credential.
err = hankoWebAuthn.DeleteCredential(credentialId)

// Search for credentials.
query = webauthn.NewCredentialQuery().
    WithUserId(userId).
    WithPageSize(pageSize).
    WithPage(page)

credentials, err = hankoWebAuthn.ListCredentials(query)
```

### Passlink usage

The Hanko Authentication API offers Passlinks as another form passwordless authentication. Instead of using a password,
the user simply clicks on a Passlink sent in an email or a message to login.

Performing a Passlink based authentication flow involves: creation a Passlink client, initialization (i.e. 
generation of Passlink and delivery of a message containing the Passlink to the user), confirmation through the user, and
finalization of the Passlink.

#### Create a new Hanko API Passlink Client

Create an account with Hanko and log into the [Hanko Console](https://console.hanko.io). The console allows you to 
create a new relying party, which also spins up a new Hanko Authentication instance (it's free, you can just try it
out). To connect the client you need the
API URL of your instance, and it is also necessary to create a new API key in the settings of your relying party to
authenticate the client. For further information, see
[Getting started](https://docs.hanko.io/gettingstarted).

Once your instance is operational, and you obtained the API key (Key-ID and a secret) from the console, provide the
Key-ID via the [`WithHmac`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#Client.WithHmac) option
and pass the secret directly into the constructor:

```go
var apiUrl string // e.g. "https://e7cd3792-dd57-4103-a6e3-4ac66fcd00f1.authentication.hanko.io"
var secret string // e.g. "FxCcAWGTEnT2P6mFbnQW4U..."

hankoPasslink = passlink.NewClient(apiUrl, secret).
    WithHmac(hmacApiKeyId). // e.g. "66dc69ec-67dc-4901-aab3-a2bf4927671e"
    WithHttpClient(httpClient). // use a customized http client
    WithLogger(logger) // use a customized logger
```

#### Passlink initialization

Initialize the Passlink using a LinkRequest:

```go
request := &passlink.LinkRequest{
    // The ID of the user to initialize a Passlink for. Solely used as a correlation identifier, since the Hanko
    // API itself does not manage user data. Must be provided by the client (relying party).
    UserID:     "d390f01d-782c-4fea-854a-abedfd5860c0",
    // Determines the communication channel through which Passlinks are delivered to the user.
    // Currently, only `email` is supported.
    Transport:  "email",
    // The recipient address the message containing the Passlink should be sent to
    Email:      "john.doe@example.com",
    // The relying party URL to redirect to after a user has confirmed (clicked) a Passlink (see "Passlink confirmation").
    // Represents the handler implemented by the relying party that performs Passlink finalization.
    // Must be a URL that has been configured by the relying party as a valid redirect URL in the Hanko Console.
    // Can be omitted in the initialization request if configured as a default redirect URL in the Hanko Console.
    RedirectTo: "https://example.com/passlink/finalize"
  }

passlink, apiErr := hankoPasslink.InitializePasslink(passlinkRequest)
```

For an in-depth description of available fields on the LinkRequest, please consult the LinkRequest code documentation
or visit our [API reference](https://docs.hanko.io/api/passlink#operation/passlinkInit).

#### Passlink confirmation

Confirmation consists of the user clicking the link delivered in the message during initialization.
Clicking the link confirms the Passlink by issuing a request to the Hanko Authentication API. The Hanko API then
redirects the user to a previously configured target in the scope of the relying party 
application the user wants to authentication with. This handler finalizes
the flow (see [Passlink finalization](#passlink-finalization)) and it must be implemented by the relying party.

#### Passlink finalization

As mentioned in [Passlink confirmation](#passlink-confirmation), the relying party must provide a handler that 
finalizes the Passlink flow (see also the `RedirectTo` attribute of the initialization request). The Hanko API appends the 
ID of the Passlink to finalize to the redirect URL after confirmation. Extract it and use it to finalize the Passlink 
flow with the Hanko API:

```go
passlink, apiErr := hankoPasslink.FinalizePasslink(linkId)
```

For a more complete implementation guide, please see the [Hanko Docs](https://docs.hanko.io/passlink/implementation).

## Examples

### WebAuthn examples

For demonstration purposes we're using the web application framework [Gin](https://github.com/gin-gonic/gin), just to
give you an idea how you can integrate the Hanko API Client. If you are interested in a full working example, check out
the [Quick Start App](https://github.com/teamhanko/hanko-webauthn-quickstart-golang).

In these examples, we will develop a small HTTP API that will be able to register and authenticate  with
WebAuthn credentials using just a few lines of code. We will see how to communicate
with the [Hanko Authentication API](https://docs.hanko.io/overview) and exchange data with the
[WebAuthn Authentication API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API).

The comments in the code samples below often mention options you can use, they are not explained in this place.
Please see the [code documentation](https://pkg.go.dev/github.com/teamhanko/hanko-go) for more details. If you
are not yet familiar with FIDO2/WebAuthn, be sure to visit [Hanko Docs](https://docs.hanko.io) first.

Before we get started, [check](https://caniuse.com/#feat=webauthn) that you are using a FIDO-compatible browser and that you have an authentication device 
available. This can be either a roaming authenticator, a platform authenticator, or a virtual 
one (e.g. [Chrome DevTools](https://developers.google.com/web/tools/chrome-devtools/webauthn)).

In addition to the [installation steps](#installation), add Gin to your dependencies:
```shell
$ go get -u github.com/gin-gonic/gin
```

Let's go ahead by creating a basic skeleton for our application:

```go
package example

import (
    "github.com/gin-gonic/gin"
    "github.com/teamhanko/hanko-go/webauthn"
)

func main() {
    var apiUrl string // e.g. "https://e7cd3792-dd57-4103-a6e3-4ac66fcd00f1.authentication.hanko.io"
    var secret string // e.g. "FxCcAWGTEnT2P6mFbnQW4U..."
    var keyId string // e.g. "65a3eba6-22cb-4c35-9881-b21fac6acfd0"
	
    hanko := webauthn.NewClient(apiUrl, secret).WithHmac(keyId)
    // Client options:
    //  - WithHmac(hmacApiKeyId)
    //  - WithHttpClient(httpClient)
    //  - WithLogger(logger)
    //  - WithoutLogs()
    //  - WithLogLevel(level)
    //  - WithLogFormatter(formatter)
	
    r := gin.Default()
    r.POST("/registration_initialize", func(c *gin.Context) { /* ... */ }) 
    // ...
}
```

If you are wondering how to get the secret to create the API client, please read the relevant 
[section](#create-a-new-hanko-api-client).

#### Example of how to register credentials

You may set up a route to initialize the credential registration that first creates a
[`RegistrationInitializationRequest`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#RegistrationInitializationRequest) 
that contains the user information, as well as information about the authenticator
to be used. Issue this request to the Hanko Authentication API using the Hanko API Client and pass the result to the
browser:

```go 
r.POST("/registration_initialize", func(c *gin.Context) {
    // To create the user object you'll need a userId and a userName. The userName usually 
    // comes either from a form a user provides when registering for the first time, or from your existing session 
    // store or database, as well as a related userId, which may needs to be generated if it is a new user. 
    userId := "65a3eba6-22cb-4c35-9881-b21fac6acfd0"
    userName := "alice@example.com"
    
    user := webauthn.NewRegistrationInitializationUser(userId, userName)
    // RegistrationInitializationUser options:
    //  - WithDisplayName(displayName)
    
    // Optionally refine which authenticator devices are allowed.
    // authenticatorSelection := webauthn.NewAuthenticatorSelection()
    // AuthenticatorSelection options:
    //  - WithUserVerification(userVerification)
    //  - WithAuthenticatorAttachment(authenticatorAttachment)
    //  - WithRequireResidentKey(requireResidentKey)
    
    request := webauthn.NewRegistrationInitializationRequest(user)
    // RegistrationInitializationRequest options:
    //  - WithAuthenticatorSelection(authenticatorSelection)
    //  - WithConveyancePreference(conveyancePreference)
    
    // Send the request to the Hanko Authentication API.
    response, apiErr := hanko.InitializeRegistration(request)
    
    if apiErr != nil {
        c.Code(apiErr.StatusCode) // pass the status code
        return
    }
    
    c.JSON(200, response) // ok
})
```

The [`RegistrationInitializationResponse`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#RegistrationInitializationResponse)
returned by `/registration_initialize`, represents [`PublicKeyCredentialCreationOptions`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredentialCreationOptions)
to be used to sign the challenge via the [`navigator.credentials.create()`](https://developer.mozilla.org/en-US/docs/Web/API/CredentialsContainer/create) browser function. 
A [`PublicKeyCredential`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredential) will be returned, which is
represented as a [`RegistrationFinalizationRequest`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#RegistrationFinalizationRequest).
Let´s set up a second route to verify the finalization request:

```go
r.POST("/registration_finalize", func (c *gin.Context) {
    request, err := webauthn.ParseRegistrationFinalizationRequest(c.Request.Body)
    if err != nil {
        c.Code(401) // bad request
        return
    }
    
    response, apiErr := hanko.FinalizeRegistration(request)
    if err != nil {
        c.Code(apiErr.StatusCode)
        return
    }
    
    c.JSON(200, response)
})
```

Almost done. You'll need some client code to bring everything together. In principle, you pass the request we retrieve 
from the `/registration_initialize` route to the WebAuthn Authentication API, the user is prompted for 
the gesture, then we send the resulting [`PublicKeyCredential`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredential) 
to `/registration_finalize` for verification.

Please note that the WebAuthn Authentication API requires data that looks like JSON but contains binary data, represented as
ArrayBuffers that needs to be encoded. So we can't pass the [`RegistrationInitializationResponse`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#RegistrationInitializationResponse)
directly as [`PublicKeyCredentialCreationOptions`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredentialCreationOptions),
but you can use the [Hanko WebAuthn Library](https://github.com/teamhanko/hanko-webauthn) 
that wraps the WebAuthn Authentication API and encodes / decodes the data.

Install the library:
```shell
$ npm install --save @teamhanko/hanko-webauthn
```

Now we can apply the following code to perform the registration:

```javascript
import {create} from "@teamhanko/hanko-webauthn"

// register() sends a request to "/registration_initialize" and passes the response to the create() function.
// The user will be promted to perform a gesture. Finally sends the resulting publicKeyCredential to 
// "/registration_finalize" for verification. This example does not handle errors.
async register() {
    const request = { body: JSON.stringify({ userName: "alice@example.com",... }), method: "POST",... }
    const registrationRequest = fetch("/registration_initialize", request)
    const creationOptions = await registrationRequest.json();
    const authenticatorResponse = await create(creationOptions);
    await fetch("/registration_finalize", { body: JSON.stringify(authenticatorResponse), method: "POST",... })
}
```

At this point you might try calling the [`ListCredentials`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#Client.ListCredentials)
function to see if a new credential appears in the list.

#### Example of how to handle the authentication

The authentication works in the same way as the registration. You'll need two routes to initialize and finalize the
authentication, and they are also used to satisfy the WebAuthn Authentication API.

With the [`AuthenticationInitializationRequest`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#AuthenticationInitializationRequest)
the user is optional. It does not need to be specified if a resident key is to be used. Otherwise, determine the
`userId` and attach the user to the request: 

```go
r.POST("/authentication_initialize", func(c *gin.Context) {
    // During authentication, the userName is usually entered by the user and sent in the request. Look up the 
    // related userId. 
    userId := "65a3eba6-22cb-4c35-9881-b21fac6acfd0"
    user := webauthn.NewAuthenticationInitializationUser(userId)

    request := webauthn.NewAuthenticationInitializationRequest().WithUser(user)
    // AuthenticationInitializationRequest options:
    //  - WithUser(user)
    //  - WithUserVerification(userVerification)
    //  - WithAuthenticatorAttachment(authenticatorAttachment)

    response, apiErr := hanko.InitializeAuthentication(request)
    if apiErr != nil {
        c.Code(apiErr.StatusCode)
        return
    }

    c.JSON(200, response)
})
```

The code to finalize the authentication is quite similar to finalizing a registration.

The [`AuthenticationInitializationResponse`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#AuthenticationFinalizationResponse)
returned by `/authentication_initialize`, represents [`PublicKeyCredentialRequestOptions`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredentialRequestOptions)
to be used to sign the challenge via the [`navigator.credentials.get()`](https://developer.mozilla.org/en-US/docs/Web/API/CredentialsContainer/get)
browser function. A [`PublicKeyCredential`](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredential) 
will be returned, which is represented as an [`AuthenticationInitializationRequest`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#AuthenticationInitializationRequest).
Let´s set up a new route to verify the finalization request:

```go
r.POST("/authentication_finalize", func(c *gin.Context) {
    request, err := webauthn.ParseAuthenticationFinalizationRequest(c.Request.Body)
    if err != nil {
        c.Code(401)
        return
    }

    response, apiErr := hanko.FinalizeAuthentication(request)
    if apiErr != nil {
        c.Code(apiErr.StatusCode)
        return
    }

    c.JSON(200, response)
})
```

The client code looks almost the same as in the case of a registration:

```javascript
import {get} from "@teamhanko/hanko-webauthn"

// authenticate() sends a request to "/authentication_initialize" and passes the response to the get() function.
// The user will be promted to perform a gesture. Finally sends the resulting credential to 
// "/authentication_finalize" for verification. This example does not handle errors.
async authenticate() {
    const request = { body: JSON.stringify({ userName: "alice@example.com",... }), method: "POST",... }
    const authenticationRequest = fetch("/authentication_initialize", request)
    const requestOptions = await authenticationRequest.json();
    const authenticatorResponse = await get(requestOptions);
    await fetch("/authentication_finalize", { body: JSON.stringify(authenticatorResponse), method: "POST",... })
}
```

If it works correctly, you should now be able to authenticate with your already registered credential. After the 
authentication was successful the LastUsed value of the credential should change. To check this, call [`GetCredential`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#Client.GetCredential) 
or [`ListCredentials`](https://pkg.go.dev/github.com/teamhanko/hanko-go/webauthn#Client.ListCredentials).

To monitor the usage of your relying party you can use the dashboard in the [Hanko Console](https://console.hanko.io).

### Passlink examples

For an in-depth Passlink example, please see the implementation guide in the [Hanko Docs](https://docs.hanko.io/passlink/implementation).

## Support

If you need help, have any questions, or have noticed an issue, please do not hesitate to
[email](mailto:support@hanko.io) us.


