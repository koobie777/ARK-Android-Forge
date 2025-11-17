//go:build !windows

package preflight

import "syscall"

func readUlimit() (uint64, error) {
	var r syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &r); err != nil {
		return 0, err
	}
	return r.Cur, nil
}
