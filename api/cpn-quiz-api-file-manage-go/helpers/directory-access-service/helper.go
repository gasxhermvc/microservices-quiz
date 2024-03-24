package directoryaccess

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func protectPath(path string) string {
	xpath := prettyPath(path)

	//=>Protection path
	if !strings.HasPrefix(xpath, "./") &&
		!strings.Contains(xpath, "://") &&
		!strings.Contains(xpath, ":/") &&
		!strings.Contains(xpath, ".:\\") {
		xpath = "./" + xpath
	}

	if strings.HasSuffix(xpath, "/") {
		xpath = strings.TrimSuffix(xpath, "/")
	}

	return strings.Trim(xpath, " ")
}

func prettyPath(path string) string {
	return strings.Replace(strings.Replace(path, "\\\\", "/", -1), "\\", "/", -1)
}

func cleansingPath(cleanPath string, path string) string {
	return strings.Replace(path, cleanPath, "", -1)
}

func generateRenameFile() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%s-%d", now.Format("20060102150405"), randInt(101, 999))
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
