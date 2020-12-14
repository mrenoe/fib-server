package main

import "testing"

func Test_solveFib(t *testing.T) {
	type args struct {
		n            int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Nominal Test",
			args: args{0},
			want: "0", wantErr: false},
		{name: "Fib at 95",
			args: args{95},
			want: "31940434634990099905", wantErr: false},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solveFib(tt.args.n)
			if got != tt.want {
				t.Errorf("findCustomerMarket() = %v, want %v", got, tt.want)
			}
		})
	}
}