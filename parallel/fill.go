package parallel

import (
	"fmt"
	"strings"
	"sync"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/mvm/progress"
)

// Fillable types have a Fill method which downloads extra data for the object
// from online sources
type Fillable interface {
	Fill(settings *htmlparsing.Settings) error
}

// MapFill calls Fill on all objects passed in the in channel, making at most
// maxRequests parallel calls.
// It increments progress each time it finishes a single call, and calls
// Done in the end.
func MapFill(
	in <-chan Fillable,
	maxRequests int,
	progress progress.Progress,
	settings *htmlparsing.Settings,
) error {
	defer progress.Done()

	out := make(chan Fillable)
	errorStream := make(chan error)

	parallelRequests := sync.WaitGroup{}
	parallelRequests.Add(maxRequests)

	for i := 0; i < maxRequests; i++ {
		go func() {
			defer parallelRequests.Done()

			for item := range in {
				err := item.Fill(settings)
				if err != nil {
					errorStream <- err
				}
				out <- item
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		parallelRequests.Wait()
		close(done)
	}()

	var errors []string

	for {
		select {
		case <-out:
			progress.Add(1)
		case err := <-errorStream:
			errors = append(errors, err.Error())
		case <-done:
			if len(errors) > 0 {
				return fmt.Errorf("%s", strings.Join(errors, ", "))
			}
			return nil
		}
	}
}
