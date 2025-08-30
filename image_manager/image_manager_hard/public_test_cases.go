package main

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Вызов",
		check: func() bool {
			return true
		},
	},
}
