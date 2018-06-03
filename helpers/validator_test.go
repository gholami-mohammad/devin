package helpers

import (
	"testing"
)

func TestValidator_IsValidUsernameFormat(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		v    Validator
		args args
		want bool
	}{
		{
			name: "ok username",
			args: args{username: "mgh"},
			want: true,
		},
		{
			name: "invalid uppercase",
			args: args{username: "UPPERCASE"},
			want: false,
		},
		{
			name: "invalid chars",
			args: args{username: "mgh$5"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{}
			if got := v.IsValidUsernameFormat(tt.args.username); got != tt.want {
				t.Errorf("Validator.IsValidUsernameFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsValidEmailFormat(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		v    Validator
		args args
		want bool
	}{
		{
			name: "OK",
			args: args{email: "info@example.com"},
			want: true,
		},
		{
			name: "empty",
			args: args{email: ""},
			want: false,
		},
		{
			name: "invalid email",
			args: args{email: "info"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{}
			if got := v.IsValidEmailFormat(tt.args.email); got != tt.want {
				t.Errorf("Validator.IsValidEmailFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
