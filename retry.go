package daap

import "fmt"

// `func retry` wraps the function given as the first parameter to retry it until the count exeecds
// its "RetryCount". The "RetryCount" is 0 unless it is specified intentionally.
func (c *Container) retry(do func() error, count int, lasterror error) error {
	if count > c.RetryCount {
		return fmt.Errorf("retry count exceeded (%d) with last error: %v", c.RetryCount, lasterror)
	}
	// Try execute given function.
	err := do()
	if err == nil {
		return nil
	}
	// If this container does NOT claim retrying, return the lasterror straightly.
	if c.RetryCount == 0 {
		return err
	}
	return c.retry(do, count+1, err)
}
