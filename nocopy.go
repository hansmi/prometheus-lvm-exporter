package main

import "sync"

// https://stackoverflow.com/a/52495303
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

var _ sync.Locker = (*noCopy)(nil)
