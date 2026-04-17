package uuid_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/WarnetBes/cursor-tool/internal/uuid"
)

var uuidV4Re = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func TestGenerateUUID_Format(t *testing.T) {
	id, err := uuid.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if !uuidV4Re.MatchString(id) {
		t.Errorf("UUID v4 format mismatch: %q", id)
	}
}

func TestGenerateUUID_Uniqueness(t *testing.T) {
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		id, err := uuid.Generate()
		if err != nil {
			t.Fatalf("Generate() error at %d: %v", i, err)
		}
		if seen[id] {
			t.Fatalf("UUID collision at iteration %d: %q", i, id)
		}
		seen[id] = true
	}
}

func TestGenerateUUID_Version4Bits(t *testing.T) {
	for i := 0; i < 100; i++ {
		id, _ := uuid.Generate()
		parts := strings.Split(id, "-")
		if parts[2][0] != '4' {
			t.Errorf("version nibble not '4' in: %q", id)
		}
		v := parts[3][0]
		if v != '8' && v != '9' && v != 'a' && v != 'b' {
			t.Errorf("variant nibble not RFC4122 in: %q", id)
		}
	}
}
