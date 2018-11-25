package main

// foo
func introduceResultBasic() introduceResultBasicResult {
	return introduceResultBasicResult{"a", "b"}
}

type introduceResultBasicResult struct {
	Field0 string
	Field1 string
}
