package recommend

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const asyFolderName = ".asy"

// RootPath is the root path of the asy(apollo synchronizer). it panics if
// the user home is not exist.
func RootPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(home, asyFolderName)
}

// GenerateAppPath generates the path of the app.
func GenerateAppPath(root string, appId string, env, cluster string) string {
	return fmt.Sprintf("%s/%s", root, strings.Join([]string{appId, env, cluster}, "-"))
}
