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
		
	}{
		{name: "Nominal Test",
			args: args{0},
			want: "0"},
		{name: "Fib(95)",
			args: args{95},
			want: "31940434634990099905"},
		{name: "Fib(96)",
			args: args{96},
			want: "51680708854858323072"},
		{name: "Additive Cache",
			args: args{97},
			want: "83621143489848422977"},
		{name: "Fib(95) cache",
			args: args{95},
			want: "31940434634990099905"},

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