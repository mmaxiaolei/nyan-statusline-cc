package formatter

import "testing"

func TestFormatCost(t *testing.T) {
	cases := []struct {
		name string
		cost float64
		want string
	}{
		{"zero", 0, "$0.0000"},
		{"negative", -5.0, "$0.0000"},
		{"tiny", 0.001, "$0.0010"},
		{"small", 0.005, "$0.0050"},
		{"below_one_cent", 0.0099, "$0.0099"},
		{"one_cent", 0.01, "$0.010"},
		{"mid_range", 0.123, "$0.123"},
		{"just_below_one", 0.999, "$0.999"},
		{"one_dollar", 1.0, "$1.00"},
		{"large", 99.99, "$99.99"},
		{"very_large", 12345.678, "$12345.68"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatCost(tc.cost)
			if got != tc.want {
				t.Errorf("FormatCost(%f) = %q, want %q", tc.cost, got, tc.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		name string
		ms   int64
		want string
	}{
		{"zero", 0, "0s"},
		{"negative", -1000, "0s"},
		{"sub_second", 500, "0s"},
		{"one_second", 1000, "1s"},
		{"seconds", 45000, "45s"},
		{"one_minute", 60000, "1m0s"},
		{"minutes_and_seconds", 150000, "2m30s"},
		{"just_under_hour", 3599000, "59m59s"},
		{"one_hour", 3600000, "1h0m"},
		{"hours_and_minutes", 5430000, "1h30m"},
		{"large", 86400000, "24h0m"},
		{"very_large", 360000000, "100h0m"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatDuration(tc.ms)
			if got != tc.want {
				t.Errorf("FormatDuration(%d) = %q, want %q", tc.ms, got, tc.want)
			}
		})
	}
}

func TestFormatTokens(t *testing.T) {
	cases := []struct {
		name   string
		tokens int64
		want   string
	}{
		{"zero", 0, "0"},
		{"negative", -100, "0"},
		{"small", 1, "1"},
		{"hundreds", 999, "999"},
		{"one_k", 1000, "1.0k"},
		{"mid_k", 5500, "5.5k"},
		{"near_10k", 9999, "10.0k"},
		{"ten_k", 10000, "10k"},
		{"large", 150000, "150k"},
		{"million", 1000000, "1000k"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatTokens(tc.tokens)
			if got != tc.want {
				t.Errorf("FormatTokens(%d) = %q, want %q", tc.tokens, got, tc.want)
			}
		})
	}
}
