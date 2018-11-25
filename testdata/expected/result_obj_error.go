package main

import "errors"

// foo
func introduceResultWithError() (introduceResultWithErrorResult, error) {
	var someCondition = true
	if someCondition {
		return introduceResultWithErrorResult{"a", "b"}, errors.New("some error")
	}
	return introduceResultWithErrorResult{"a", "b"}, nil
}

type introduceResultWithErrorResult struct {
	Field0 string
	Field1 string
}
