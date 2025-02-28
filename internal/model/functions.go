package model

import "time"

func Debounce(f func(), delay time.Duration) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, f)
	}
}

func Debounce2[T any, K any](f func(one T, second K), delay time.Duration) func(one T, second K) {
	var timer *time.Timer
	return func(one T, second K) {
		if timer != nil {
			timer.Stop()
		}
		var arg = func() {
			f(one, second)
		}
		timer = time.AfterFunc(delay, arg)
	}
}
