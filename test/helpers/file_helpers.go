package helpers

import (
	"os"
)

func CreateTempTextFile(
	content string,
) (*os.File) {
	tempFile, _ := os.CreateTemp("", "testfile-*.txt")
	defer tempFile.Close()

	tempFile.WriteString(content)

	return tempFile
}