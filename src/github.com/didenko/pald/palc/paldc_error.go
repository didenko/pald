package palc

import "fmt"

type PalcError struct {
	message string
	errCode int
	svrName string
	svrPort uint16
}

func (pe *PalcError) Error() string {
	return fmt.Sprintf(
		"server: %q, port: %d, error: %q",
		pe.svrName,
		pe.svrPort,
		pe.message)
}

func newFromResp(svr string, port uint16, status int, body string) error {

	if status < 300 {
		return nil
	}

	return &PalcError{body, status, svr, port}
}

func newFromError(svr string, port uint16, err error) error {

	if err == nil {
		return nil
	}

	return &PalcError{err.Error(), 1, svr, port}
}
