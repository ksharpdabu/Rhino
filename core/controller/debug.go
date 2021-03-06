// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/rhino/core/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Debug controller
func Debug(c *gin.Context) {
	var bodyBytes []byte

	// Workaround for issue https://github.com/gin-gonic/gin/issues/1651
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	parameters := make(map[string]string)
	headers := make(map[string]string)

	for k, v := range c.Request.URL.Query() {
		parameters[k] = strings.Join(v, ", ")
	}

	for k, v := range c.Request.Header {
		headers[k] = strings.Join(v, ", ")
	}

	route := model.GetRoute(c.FullPath(), "", parameters)

	rand.Seed(time.Now().UnixNano())

	failCount, _ := strconv.Atoi(strings.Replace(route.Chaos.FailRate, "%", "", -1))

	if rand.Intn(100) < failCount {
		log.WithFields(log.Fields{
			"method":     c.Request.Method,
			"url":        c.Request.URL.Path,
			"headers":    headers,
			"parameters": parameters,
			"body":       string(bodyBytes),
		}).Info("Failed Request")

		c.Status(http.StatusInternalServerError)
		return
	}

	latencySeconds, _ := strconv.Atoi(strings.Replace(route.Chaos.Latency, "s", "", -1))

	time.Sleep(time.Duration(latencySeconds) * time.Second)

	log.WithFields(log.Fields{
		"method":     c.Request.Method,
		"url":        c.Request.URL.Path,
		"headers":    headers,
		"parameters": parameters,
		"body":       string(bodyBytes),
	}).Info("Request Success")

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
