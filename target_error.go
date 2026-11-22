package errlay

import "testing"

func TestTarget(t *testing.T) {
	layer := "database"
	emptyLayer := ""
	code := uint8(42)
	zeroCode := uint8(0)

	tests := []struct {
		name    string
		layer   *string
		code    *uint8
		wantErr bool
	}{
		{name: "layer only", layer: &layer},
		{name: "code only", code: &code},
		{name: "both fields", layer: &layer, code: &code},
		{name: "empty layer is specified", layer: &emptyLayer},
		{name: "zero code is specified", code: &zeroCode},
		{name: "neither field", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			layer, code := tt.layer, tt.code

			// Act
			target, err := Target(layer, code)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatal("Target() error = nil, want an error")
				}
				if target != nil {
					t.Fatalf("Target() target = %v, want nil", target)
				}
				return
			}
			if err != nil {
				t.Fatalf("Target() unexpected error: %v", err)
			}
			got, ok := target.(*errTarget)
			if !ok {
				t.Fatalf("Target() returned %T, want *errTarget", target)
			}
			if got.layer != layer || got.code != code {
				t.Errorf("Target() fields = (%v, %v), want (%v, %v)", got.layer, got.code, layer, code)
			}
		})
	}
}

func TestErrTargetError(t *testing.T) {
	tests := []struct {
		name string
		obj  *errTarget
		want string
	}{
		{name: "empty representation", obj: &errTarget{}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			target := tt.obj

			// Act
			got := target.Error()

			// Assert
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
