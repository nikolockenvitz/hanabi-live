// Miscellaneous subroutines

package main

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strings"

	sentry "github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func executeScript(scriptName string) error {
	scriptPath := path.Join(projectPath, scriptName)
	cmd := exec.Command(scriptPath)
	cmd.Dir = projectPath
	outputBytes, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(outputBytes))
	hLog.Info(
		fmt.Sprintf("Script \"%v\" completed.", scriptName),
		zap.String("output", output),
	)
	if err != nil {
		// The "cmd.CombinedOutput()" function will throw an error if the return code is not equal
		// to 0
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("scriptOutput", output)
		})
		return fmt.Errorf("failed to execute \"%v\": %w", scriptPath, err)
	}
	return nil
}

// From: https://stackoverflow.com/questions/38554353/how-to-check-if-a-string-only-contains-alphabetic-characters-in-go
var isAlphanumericHyphen = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString

// From: https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// From: https://mrekucci.blogspot.com/2015/07/dont-abuse-mathmax-mathmin.html
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}