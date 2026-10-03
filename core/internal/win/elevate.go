package win


import "golang.org/x/sys/windows"

func IsElevated() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}