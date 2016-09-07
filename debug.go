package htmlparsing

import (
	"fmt"
	"os"

	"github.com/jbowtie/gokogiri/xml"
)

func DumpHTML(node xml.Node, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = fmt.Fprintf(f, node.InnerHtml())
	if err != nil {
		panic(err)
	}
}
