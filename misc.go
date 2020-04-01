package cynhid

import "C"
import (
	"fmt"
)

// Error represents a generic error
type Error C.int

func (err Error) Error() string {
	return fmt.Sprintf("hidapi: %s", err)
}

func errFromErrno(errno C.int) error {
	err := Error(errno)
	if errno == 0 {
		return nil
	}
	return err
}
