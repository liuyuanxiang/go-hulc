package util

import "fmt"

func GetPortString(port int64) string {
	return fmt.Sprintf(":%d", port)
}
