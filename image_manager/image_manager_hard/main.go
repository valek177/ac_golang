package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
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
