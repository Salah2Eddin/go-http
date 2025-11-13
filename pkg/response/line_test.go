package response

import "testing"

func TestNewStatusLine(t *testing.T) {
	tests := []struct {
		code     int
		wantCode int
		wantVer  string
		wantText string
	}{
		{200, 200, "HTTP/1.1", "OK"},
		{404, 404, "HTTP/1.1", "Not Found"},
		{999, 999, "HTTP/1.1", ""}, // unknown code
	}

	for _, tt := range tests {
		sl := NewStatusLine(tt.code)
		if sl.Code != tt.wantCode {
			t.Errorf("Code = %d, want %d", sl.Code, tt.wantCode)
		}
		if sl.Version != tt.wantVer {
			t.Errorf("Version = %q, want %q", sl.Version, tt.wantVer)
		}
		if sl.Phrase != tt.wantText {
			t.Errorf("Phrase = %q, want %q", sl.Phrase, tt.wantText)
		}
	}
}
