package metrics

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ResponseDuration struct {
	Path string `json:"path"`
	Time int64  `json:"average_time"`
}

type ErrorCount struct {
	Status int `json:"status_code"`
	Count  int `json:"count"`
}

type Metrics struct {
	ResponseDuration []ResponseDuration `json:"response_metrics"`
	ErrorCount       []ErrorCount       `json:"error_counts"`
}

var responeMetrics = make(map[string][]int64)
var errorMetrics = make(map[int]int)

func ResponseMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		path := c.FullPath()
		if _, ok := responeMetrics[path]; !ok {
			responeMetrics[path] = make([]int64, 0)
		}
		list := responeMetrics[path]
		list = append(list, time.Since(now).Nanoseconds())
		responeMetrics[path] = list
	}
}

func CaptureErrorMetrics(statusCode int) {
	if _, ok := errorMetrics[statusCode]; !ok {
		errorMetrics[statusCode] = 0
	}
	errorMetrics[statusCode] += 1
}

func GetResponseMetrics(c *gin.Context) {
	responseDuration := make([]ResponseDuration, 0)
	for k, v := range responeMetrics {
		total := int64(0)
		for _, v := range v {
			total += v
		}
		duration := ResponseDuration{
			Path: k,
			Time: total / int64(len(v)),
		}
		responseDuration = append(responseDuration, duration)
	}
	errorCounts := make([]ErrorCount, 0)
	for k, v := range errorMetrics {
		errorCount := ErrorCount{
			Status: k,
			Count:  v,
		}
		errorCounts = append(errorCounts, errorCount)
	}

	c.JSON(http.StatusOK, gin.H{"metrics": Metrics{
		ResponseDuration: responseDuration,
		ErrorCount:       errorCounts,
	}})
}
