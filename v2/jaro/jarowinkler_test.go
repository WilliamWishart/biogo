package jaro

import "testing"

func TestJaroSimilarity_Basic(t *testing.T) {
	cases := []struct {
		a, b   string
		expect float32
	}{
		{"", "", 0.0},
		{"abc", "abc", 1.0},
		{"abc", "acb", 0.5555556}, // corrected expected value
		{"abc", "xyz", 0.0},
		{"MARTHA", "MARHTA", 0.9444444},
		{"DWAYNE", "DUANE", 0.8222222},
		{"DIXON", "DICKSONX", 0.7666667},
	}
	for _, c := range cases {
		got := JaroSimilarity(c.a, c.b)
		if (c.expect == 1.0 && got != 1.0) || (c.expect == 0.0 && got != 0.0) {
			t.Errorf("JaroSimilarity(%q, %q) = %v, want %v", c.a, c.b, got, c.expect)
		} else if c.expect != 1.0 && c.expect != 0.0 {
			if diff := got - c.expect; diff > 0.01 || diff < -0.01 {
				t.Errorf("JaroSimilarity(%q, %q) = %v, want %v", c.a, c.b, got, c.expect)
			}
		}
	}
}

func TestJaroSimilarity_Unicode(t *testing.T) {
	got := JaroSimilarity("你好", "你号")
	if got <= 0.0 || got >= 1.0 {
		t.Errorf("JaroSimilarity for unicode should be between 0 and 1, got %v", got)
	}
}
