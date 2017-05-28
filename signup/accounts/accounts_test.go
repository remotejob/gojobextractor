package accounts

import "testing"

func TestGetCsv(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"test0", args{"/home/juno/gowork/src/github.com/remotejob/gojobextractor/accounts.csv"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCsv(tt.args.file)
		})
	}
}
