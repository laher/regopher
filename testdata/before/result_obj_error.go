package main

import "errors"

// foo
func introduceResultWithError() (string, string, error) {
	var someCondition = true
	if someCondition {
		return "a", "b", errors.New("some error")
	}
	return "a", "b", nil
}
