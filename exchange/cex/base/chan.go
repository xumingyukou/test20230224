package base

import (
	"clients/logger"
	"fmt"
	"strings"
)

func SendChan[T any](ch chan T, res T, title ...string) {
	content := strings.Join(title, " ")
	if len(ch) > cap(ch)/2 {
		content := fmt.Sprint(content, " channel slow:", len(ch))
		logger.Logger.Warn(content)
	}
	ch <- res
}
