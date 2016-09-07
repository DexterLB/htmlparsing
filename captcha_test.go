package htmlparsing

import (
	"os"
	"testing"
)

func TestBreakSimpleCaptcha(t *testing.T) {
	f, err := os.Open("fixtures/captcha.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	result, err := BreakSimpleCaptcha(f)
	if err != nil {
		t.Error(err)
	}

	if result != "9179" {
		t.Errorf("%s != 9179", result)
	}
}
