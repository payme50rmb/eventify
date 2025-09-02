package eventify

import (
	"reflect"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Matcher
	}{
		{
			name: "test",
			args: args{
				s: "test",
			},
			want: &Matcher{
				s: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMatcher(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatcher_Match(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		m    *Matcher
		args args
		want bool
	}{
		{
			name: "equal",
			m:    NewMatcher("foo"),
			args: args{
				target: "foo",
			},
			want: true,
		},
		{
			name: "any",
			m:    NewMatcher("*"),
			args: args{
				target: "test2",
			},
			want: true,
		},
		{
			name: "prefix-true",
			m:    NewMatcher("test*"),
			args: args{
				target: "test3",
			},
			want: true,
		},
		{
			name: "prefix-false",
			m:    NewMatcher("test*"),
			args: args{
				target: "1test3",
			},
			want: false,
		},
		{
			name: "suffix-true",
			m:    NewMatcher("*test"),
			args: args{
				target: "1test",
			},
			want: true,
		},
		{
			name: "suffix-false",
			m:    NewMatcher("*test"),
			args: args{
				target: "1test2",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Match(tt.args.target); got != tt.want {
				t.Errorf("Matcher.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
