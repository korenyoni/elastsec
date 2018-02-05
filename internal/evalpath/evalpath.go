package evalpath

import (
    "path"
    "strings"
    "fmt"
)

func Eval(cwd string, parentPath string, childPath string) string {
    filePath := parentPath
    if strings.HasPrefix(parentPath, "/") && strings.HasPrefix(childPath,parentPath) {
        filePath = childPath
    } else if strings.HasPrefix(parentPath, "/") {
        filePath = parentPath + "/" + childPath
    } else {
        filePath = cwd + "/" + childPath
    }

    return path.Clean(filePath)
}
