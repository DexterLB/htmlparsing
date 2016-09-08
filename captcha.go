package htmlparsing

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// BreakSimpleCaptcha uses the tesseract OCR command to recognise text in
// a simple captcha
func BreakSimpleCaptcha(image io.Reader) (string, error) {
	var textData bytes.Buffer

	ocrCommand := exec.Command("tesseract", "-", "-")
	ocrCommand.Stdin = image
	ocrCommand.Stdout = &textData

	err := ocrCommand.Run()
	if err != nil {
		return "", fmt.Errorf("Error running OCR: %s", err)
	}

	return strings.TrimSpace(textData.String()), nil
}
