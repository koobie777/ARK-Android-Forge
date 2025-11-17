//go:build windows

package preflight

func readUlimit() (uint64, error) {
	return 0, ErrUnsupported
}
