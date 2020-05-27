package param

import "testing"

func Test_parseHumanReadableSize(t *testing.T) {
	type args struct {
		humanReadable string
	}
	successTests := []struct {
		name string
		args args
		want int64
	}{
		{"zero int", args{"0"}, 0},
		{"zero float", args{"0.0"}, 0},
		{"one byte int", args{"1"}, 1},
		{"one byte float", args{"1.0"}, 1},
		{"1.4 round down", args{"1.4"}, 1},
		{"1.5 round up", args{"1.5"}, 2},
		{"1 kibibyte", args{"1K"}, 1024},
		{"1 mebibyte", args{"1M"}, 1024 * 1024},
		{"1 gibibyte", args{"1G"}, 1024 * 1024 * 1024},
		{"1 tebibyte", args{"1T"}, 1024 * 1024 * 1024 * 1024},
		{"1 pebibyte", args{"1P"}, 1024 * 1024 * 1024 * 1024 * 1024},
		{"1 exbibyte", args{"1E"}, 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
		{"half a mebibyte", args{".5M"}, 1024 * 1024 / 2},
		{"minus one", args{"-1.0"}, -1},
		{"minus one mebibyte", args{"-1M"}, -1024 * 1024},
		// max int64 is 9223372036854775807 but the float64 truncates
		{"maximum int64", args{"9223372036849999872"}, 9223372036849999872},
	}
	failureTests := []struct {
		name string
		args args
	}{
		{"empty string", args{""}},
		{"all text", args{"abc"}},
		{"text prefix number suffix", args{"abc9"}},
		{"text prefix kibibyte suffix", args{"abcK"}},
		{"overflow", args{"9223372036854775807"}},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := parseHumanReadableSize(tt.args.humanReadable); err != nil {
				t.Errorf("parseHumanReadableSize() error = %v", err)
			} else if got != tt.want {
				t.Errorf("parseHumanReadableSize() got = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := parseHumanReadableSize(tt.args.humanReadable); err == nil {
				t.Errorf("parseHumanReadableSize() got = %v, but want an error", got)
			}
		})
	}
}
