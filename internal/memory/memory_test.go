package memory_test

import (
	"fmt"
	"strings"
	"testing"
	"wrenchi/internal/memory"
)

func TestParseMem(t *testing.T) {

	var (
		memTotal     uint64 = 16384256
		memFree      uint64 = 1234567
		memAvailable uint64 = 2345678
		buffers      uint64 = 345678
		cached       uint64 = 4567890
		swapTotal    uint64 = 987654
		swapFree     uint64 = 876543
	)

	mockData := fmt.Sprintf(
		"MemTotal:        %v kB\n"+
			"MemFree:         %v kB\n"+
			"Buffers:         %v kB\n"+
			"MemAvailable:    %v kB\n"+
			"Cached:          %v kB\n"+
			"SwapTotal:       %v kB\n"+
			"SwapFree:        %v kB\n",
		memTotal,
		memFree,
		buffers,
		memAvailable,
		cached,
		swapTotal,
		swapFree)

	reader := strings.NewReader(mockData)
	m, err := memory.ParseMemInfo(reader)
	if err != nil {
		t.Fatalf("failed to parse %s", err)
	}

	if m.Total != memTotal {
		t.Errorf("Expected %v, got %v instead.", memTotal, m.Total)
	}

	if m.Free != memFree {
		t.Errorf("Expected %v, got %v instead.", memFree, m.Free)
	}

	if m.Buffers != buffers {
		t.Errorf("Expected %v, got %v instead.", buffers, m.Buffers)
	}

	if m.Available != memAvailable {
		t.Errorf("Expected %v, got %v instead.", memAvailable, m.Available)
	}

	if m.Cached != cached {
		t.Errorf("Expected %v, got %v instead.", cached, m.Cached)
	}

	if m.SwapTotal != swapTotal {
		t.Errorf("Expected %v, got %v instead.", swapTotal, m.SwapTotal)
	}

	if m.SwapFree != swapFree {
		t.Errorf("Expected %v, got %v instead.", swapFree, m.SwapFree)
	}
}

func TestParseIds(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		want      memory.IDs
		expectErr bool
	}{
		{
			name:  "Valid ids",
			input: []string{"1000", "1000", "1000", "1000"},
			want:  memory.IDs{Real: 1000, Effective: 1000, Saved: 1000, FS: 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := memory.ParseId(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}
