# Golim - A Gin-gonic middleware for limit request amount coming from an IP

Golim is a different version of limiting the amount of requests comming to a server which is limit the number of requests coming from an IP. There are a lot of way to block spam, such as [Go Rate](https://pkg.go.dev/golang.org/x/time/rate), [Limiter](https://pkg.go.dev/github.com/ulule/limiter/v3#section-readme), etc. Personally, I prefer to limiting number of requests of a user. That's why this middleware emerge.

## Usage
With [Go module](https://go.dev/wiki/Modules) support, we could simply add:

``` golang
import (
	"github.com/gin-gonic/gin"
	"github.com/ltphat2204/golim"
)
```

Due to developing based on Gin framework, we have to import [Gin](https://github.com/gin-gonic/gin) package also. For more detail, head over to the [original page](https://github.com/gin-gonic/gin?).

Otherwise, run the following Go command to install the `golim` package

``` bash
go get -u "github.com/ltphat2204/golim"
```

In main package, we just use this middleware:

``` golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ltphat2204/golim"
	"net/http"
)

func main() {
	r := gin.Default()

	limiter := golim.NewLimiter()
	r.Use(golim.LimitRequestOnIp(limiter))

	//Testing route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	r.Run()
}
```

## How does this middleware work?
When a request comming, this golim will check the total access times IP with the map of `{Ip: TotalAccessTimes}`. In default, one IP could make at most 15 requests per second. Then, if this IP's access times beyond the limit, it will be blocked and receive following response:

``` json
{
    "message": "Hold on before sending new request"
}
```

Otherwise, next middleware is called.

## How to customize the limit?

Since the limit is based on the struct `Limiter`, we could easily customize this by following example of 5 requests per minute:

``` golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ltphat2204/golim"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()

	limiter := golim.Limiter {
		Duration: time.Minute,
		Max: 5,
	}
	r.Use(golim.LimitRequestOnIp(&limiter))

	//Testing route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	r.Run()
}
```