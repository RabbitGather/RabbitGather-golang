// lib package is a package that contains some useful tools only for this project.
package lib

func Pipe[T any](inputChan <-chan T, errChan <-chan error, f func(input T) bool) error {
	for {
		select {
		case result := <-inputChan:
			if !f(result) {
				return nil
			}
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
}
