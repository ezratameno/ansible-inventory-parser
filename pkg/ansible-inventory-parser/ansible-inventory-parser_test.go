package ansibleinventoryparser

import "testing"

func TestIsKeyVal(t *testing.T) {
	type args struct {
		row string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test key val",
			args: args{
				row: "ezra: 120.5",
			},
			want: true,
		},
		{
			name: "test map",
			args: args{
				row: "ezra:",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKeyVal(tt.args.row); got != tt.want {
				t.Errorf("IsKeyVal() = %v, want %v", got, tt.want)
			}
		})
	}
}
