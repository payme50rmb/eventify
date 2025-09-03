package eventify

import "testing"

func TestNewMatcher(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		tests   []struct {
			input string
			want  bool
		}
	}{
		{
			name:    "exact match",
			pattern: "test",
			tests: []struct {
				input string
				want  bool
			}{
				{input: "test", want: true},
				{input: "test1", want: false},
				{input: "test ", want: false},
				{input: " test", want: false},
			},
		},
		{
			name:    "prefix match",
			pattern: "test*",
			tests: []struct {
				input string
				want  bool
			}{
				{input: "test", want: true},
				{input: "test123", want: true},
				{input: "test 123", want: true},
				{input: " test", want: false},
				{input: "atest", want: false},
			},
		},
		{
			name:    "suffix match",
			pattern: "*test",
			tests: []struct {
				input string
				want  bool
			}{
				{input: "test", want: true},
				{input: "123test", want: true},
				{input: " test", want: true},
				{input: "test ", want: false},
				{input: "test1", want: false},
			},
		},
		{
			name:    "wildcard match",
			pattern: "*",
			tests: []struct {
				input string
				want  bool
			}{
				{input: "", want: true},
				{input: "test", want: true},
				{input: " ", want: true},
				{input: "*", want: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatcher(tt.pattern)
			for _, tc := range tt.tests {
				t.Run(tc.input, func(t *testing.T) {
					if got := matcher.Match(tc.input); got != tc.want {
						t.Errorf("Match(%q) = %v, want %v", tc.input, got, tc.want)
					}
				})
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
		{
			name: "contains-true",
			m:    NewMatcher("*test*"),
			args: args{
				target: "test123test",
			},
			want: true,
		},
		{
			name: "contains-false",
			m:    NewMatcher("*test*"),
			args: args{
				target: "1te st1",
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
