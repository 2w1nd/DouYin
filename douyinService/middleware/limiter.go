package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"sync"
	"time"
)

func NewLimiter(r rate.Limit, b int, t time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}
	return func(c *gin.Context) {
		// 获取限速器
		// key 除了 ip 之外也可以是其他的，例如 header，user name 等
		key := c.ClientIP()
		l, _ := limiters.LoadOrStore(key, rate.NewLimiter(r, b))
		// 这里注意不要直接使用 gin 的 context 默认是没有超时时间的
		ctx, cancel := context.WithTimeout(c, t)
		defer cancel()
		if err := l.(*rate.Limiter).Wait(ctx); err != nil {
			// 这里先不处理日志了，如果返回错误就直接 429
			log.Println("进入限流器")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, "点击太快了")
			//response.FailWithMessage("点击太快了", c)
		}
		c.Next()
	}
}
