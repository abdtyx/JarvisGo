package errdefs

import "fmt"

type ErrPermissionDenied struct {
	Sig   string
	User  int64
	Group int64
}

func (e *ErrPermissionDenied) String() string {
	return fmt.Sprintf("%v: Permission denied (user: %v, group: %v)", e.Sig, e.User, e.Group)
}
