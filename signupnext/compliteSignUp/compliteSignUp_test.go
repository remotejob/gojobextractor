package compliteSignUp

import "testing"

func TestComplite(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"test0", args{"https://stackoverflow.com/users/signup-finish?email=programming%40alastomana.fi&name=Programmer+Development&token=444d3aed-4bbc-4645-8b62-40cc88bf6864%3apmj%2fogvLaJeVVS5yb83VVgL64qg%3d&sus=head&returnurl=%2fusers%2fstory%2fcurrent&authCode=dc78b718a96d148e8dea98fa8ff3f29a5f001c15fcc5a8dc37aa01d6d555983f"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Complite(tt.args.link)
		})
	}
}
