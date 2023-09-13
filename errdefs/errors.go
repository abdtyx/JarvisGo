package errdefs

import "fmt"

type BaseErr struct {
	Sig   string
	User  int64
	Group int64
}

func (e *BaseErr) String() string {
	return fmt.Sprintf("%v (user: %v, group: %v): ", e.Sig, e.User, e.Group)
}

func (e *BaseErr) Error() string {
	return fmt.Sprintf("%v (user: %v, group: %v): ", e.Sig, e.User, e.Group)
}

type ErrPermissionDenied struct {
	BaseErr
}

func (e *ErrPermissionDenied) String() string {
	return e.String() + "Sir, you have no permission to this command."
}

func (e *ErrPermissionDenied) Error() string {
	return e.Error() + "Sir, you have no permission to this command."
}

type ErrWrongParams struct {
	BaseErr
}

func (e *ErrWrongParams) String() string {
	return e.String() + "Sir, I'm afraid that you have delivered a wrong parameter."
}

func (e *ErrWrongParams) Error() string {
	return e.Error() + "Sir, I'm afraid that you have delivered a wrong parameter."
}
