package uuid

import "testing"

func TestGenerateUuid(t *testing.T) {
	uuid := Uuid{}
	generatedUuid := uuid.GenerateUuid()

	if len(generatedUuid) == 0 {
		t.Error("Generated UUID should not be empty")
	}

	if len(generatedUuid) > 66 {
		t.Error("Generated UUID should not exceed 66 characters")
	}
}
