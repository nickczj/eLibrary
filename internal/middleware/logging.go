package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Log request details
		var bodyBytes []byte
		if c.Request.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(c.Request.Body)
			if err != nil {
				log.WithError(err).Error("Error reading request body")
			} else {
				// Reset the request body for subsequent handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		log.WithFields(log.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"headers": c.Request.Header,
			"body":    string(bodyBytes),
		}).Info("Incoming Request")

		// Capture response body
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{ResponseWriter: c.Writer, body: responseBody}
		c.Writer = writer

		// Process the request
		c.Next()

		// Log response details
		endTime := time.Since(startTime)
		log.WithFields(log.Fields{
			"status":       writer.Status(),
			"responseBody": responseBody.String(),
			"responseTime": endTime,
		}).Info("Response Details")
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)                  // Copy response to buffer
	return rw.ResponseWriter.Write(data) // Write response to client
}
