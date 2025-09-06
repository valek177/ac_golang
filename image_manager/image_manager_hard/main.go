package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() struct{} {
				return struct{}{}
			},
			func(_ struct{}) bool {
				return tt.check()
			},
		)
	}
}
