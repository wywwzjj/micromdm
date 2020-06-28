// +build !windows

package log

import "syscall"

var defaultSwapSignal = syscall.SIGUSR2
