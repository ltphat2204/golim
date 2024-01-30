package golim

import (
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
    "sync"
)

var (
	lastAccessOf map[string]time.Time
	countAccessOf map[string]int
	mutex sync.Mutex
)

//Cleaning cache after using
//Frequency: Every 5 minutes
func clearCacheAccess() {
	mutex.Lock()
    defer mutex.Unlock()

    now := time.Now()
    for key, accessTime := range lastAccessOf {
        if now.Sub(accessTime) >= 5 * time.Minute {
            delete(lastAccessOf, key)
        }
    }
}

func init() {
	lastAccessOf = make(map[string]time.Time)
	countAccessOf = make(map[string]int)
	go func() {
		for {
			clearCacheAccess()
			time.Sleep(5 * time.Minute)
		}
    }()
}

func LimitRequestOnIp(limiter *Limiter) gin.HandlerFunc {
	if limiter == nil {
		limiter = NewLimiter()
	}
	now := time.Now()

	return func(c *gin.Context) {
		// Avoid conflict when interact with map in the clear cache period
		mutex.Lock()
		defer mutex.Unlock()

		ip := c.ClientIP()
		
		if lastAccessTime, ok := lastAccessOf[ip]; ok {
			if now.Sub(lastAccessTime) < limiter.Duration {
				countAccessOf[ip] = 0
			} else {
				if countAccessOf[ip] < limiter.Max {
					countAccessOf[ip] += 1
				} else {
					c.JSON(http.StatusTooManyRequests, gin.H{"message": "Hold on before sending new request"})
					c.Abort()
				}
			}
		} else {
			lastAccessOf[ip] = now
			countAccessOf[ip] = 1
		}

		c.Next()
	}
}