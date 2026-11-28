package entity_test

import (
	"errors"
	"testing"

	"virtual-library/catalog-service/internal/domain/entity"
)

func TestNewCategory_Success(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "simple name", arg: "Fantasy", want: "Fantasy"},
		{name: "name with surrounding spaces is stored as is", arg: " X ", want: " X "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := entity.NewCategory(tt.arg)
			if err != nil {
				t.Fatalf("NewCategory() unexpected error: %v", err)
			}
			if category == nil {
				t.Fatal("NewCategory() returned nil category")
			}
			if category.ID() != 0 {
				t.Errorf("NewCategory() ID = %d, want 0", category.ID())
			}
			if category.Name() != tt.want {
				t.Errorf("NewCategory() Name = %q, want %q", category.Name(), tt.want)
			}
		})
	}
}

func TestNewCategory_ValidationError(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "empty name", arg: ""},
		{name: "spaces name", arg: "   "},
		{name: "whitespace name", arg: "\t\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := entity.NewCategory(tt.arg)
			if err == nil {
				t.Fatal("NewCategory() expected error, got nil")
			}
			if !errors.Is(err, entity.ErrValidationError) {
				t.Errorf("NewCategory() error = %v, want ErrValidationError", err)
			}
			if category == nil {
				t.Fatal("NewCategory() returned nil category, want zero-value category")
			}
			if category.ID() != 0 || category.Name() != "" {
				t.Errorf("NewCategory() category = {%d, %q}, want zero value", category.ID(), category.Name())
			}
		})
	}
}

func TestHydrateCategory(t *testing.T) {
	tests := []struct {
		name string
		id   int
		arg  string
	}{
		{name: "valid data", id: 5, arg: "Fantasy"},
		{name: "empty name is kept without validation", id: 7, arg: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := entity.HydrateCategory(tt.id, tt.arg)
			if category == nil {
				t.Fatal("HydrateCategory() returned nil")
			}
			if category.ID() != tt.id {
				t.Errorf("HydrateCategory() ID = %d, want %d", category.ID(), tt.id)
			}
			if category.Name() != tt.arg {
				t.Errorf("HydrateCategory() Name = %q, want %q", category.Name(), tt.arg)
			}
		})
	}
}

func TestApplyCategory_Success(t *testing.T) {
	t.Run("rename returns new category and keeps original", func(t *testing.T) {
		original := entity.HydrateCategory(5, "Fantasy")

		updated, err := original.ApplyCategory(entity.CategoryChanges{Name: strPtr("Sci-Fi")})
		if err != nil {
			t.Fatalf("ApplyCategory() unexpected error: %v", err)
		}
		if updated == nil {
			t.Fatal("ApplyCategory() returned nil category")
		}
		if updated.ID() != 5 {
			t.Errorf("ApplyCategory() ID = %d, want 5", updated.ID())
		}
		if updated.Name() != "Sci-Fi" {
			t.Errorf("ApplyCategory() Name = %q, want %q", updated.Name(), "Sci-Fi")
		}
		if original.Name() != "Fantasy" {
			t.Errorf("ApplyCategory() mutated original, Name = %q, want %q", original.Name(), "Fantasy")
		}
		if updated == original {
			t.Error("ApplyCategory() returned same pointer as original, want new pointer")
		}
	})

	t.Run("empty changes return copy without error", func(t *testing.T) {
		original := entity.HydrateCategory(5, "Fantasy")

		updated, err := original.ApplyCategory(entity.CategoryChanges{})
		if err != nil {
			t.Fatalf("ApplyCategory() unexpected error: %v", err)
		}
		if updated == nil {
			t.Fatal("ApplyCategory() returned nil category")
		}
		if updated.Name() != "Fantasy" {
			t.Errorf("ApplyCategory() Name = %q, want %q", updated.Name(), "Fantasy")
		}
		if updated == original {
			t.Error("ApplyCategory() returned same pointer as original, want new pointer")
		}
	})
}

func TestApplyCategory_Errors(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "empty name", arg: ""},
		{name: "spaces name", arg: "  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := entity.HydrateCategory(5, "Fantasy")

			updated, err := original.ApplyCategory(entity.CategoryChanges{Name: strPtr(tt.arg)})
			if err == nil {
				t.Fatal("ApplyCategory() expected error, got nil")
			}
			if !errors.Is(err, entity.ErrValidationError) {
				t.Errorf("ApplyCategory() error = %v, want ErrValidationError", err)
			}
			if updated != nil {
				t.Errorf("ApplyCategory() category = %+v, want nil", updated)
			}
			if original.Name() != "Fantasy" {
				t.Errorf("ApplyCategory() mutated original, Name = %q, want %q", original.Name(), "Fantasy")
			}
		})
	}
}
