package entity_test

import (
	"errors"
	"testing"

	"virtual-library/catalog-service/internal/domain/entity"
)

func TestNewAuthor_Success(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "simple name", arg: "John", want: "John"},
		{name: "name with surrounding spaces is stored as is", arg: " John ", want: " John "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author, err := entity.NewAuthor(tt.arg)
			if err != nil {
				t.Fatalf("NewAuthor() unexpected error: %v", err)
			}
			if author == nil {
				t.Fatal("NewAuthor() returned nil author")
			}
			if author.ID() != 0 {
				t.Errorf("NewAuthor() ID = %d, want 0", author.ID())
			}
			if author.Name() != tt.want {
				t.Errorf("NewAuthor() Name = %q, want %q", author.Name(), tt.want)
			}
		})
	}
}

func TestNewAuthor_ValidationError(t *testing.T) {
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
			author, err := entity.NewAuthor(tt.arg)
			if err == nil {
				t.Fatal("NewAuthor() expected error, got nil")
			}
			if !errors.Is(err, entity.ErrValidationError) {
				t.Errorf("NewAuthor() error = %v, want ErrValidationError", err)
			}
			if author == nil {
				t.Fatal("NewAuthor() returned nil author, want zero-value author")
			}
			if author.ID() != 0 || author.Name() != "" {
				t.Errorf("NewAuthor() author = {%d, %q}, want zero value", author.ID(), author.Name())
			}
		})
	}
}

func TestHydrateAuthor(t *testing.T) {
	tests := []struct {
		name string
		id   int
		arg  string
	}{
		{name: "valid data", id: 5, arg: "John"},
		{name: "empty name is kept without validation", id: 7, arg: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author := entity.HydrateAuthor(tt.id, tt.arg)
			if author == nil {
				t.Fatal("HydrateAuthor() returned nil")
			}
			if author.ID() != tt.id {
				t.Errorf("HydrateAuthor() ID = %d, want %d", author.ID(), tt.id)
			}
			if author.Name() != tt.arg {
				t.Errorf("HydrateAuthor() Name = %q, want %q", author.Name(), tt.arg)
			}
		})
	}
}

func TestApplyAuthor_Success(t *testing.T) {
	t.Run("rename returns new author and keeps original", func(t *testing.T) {
		original := entity.HydrateAuthor(5, "John")

		updated, err := original.ApplyAuthor(entity.AuthorChanges{Name: strPtr("Jane")})
		if err != nil {
			t.Fatalf("ApplyAuthor() unexpected error: %v", err)
		}
		if updated == nil {
			t.Fatal("ApplyAuthor() returned nil author")
		}
		if updated.ID() != 5 {
			t.Errorf("ApplyAuthor() ID = %d, want 5", updated.ID())
		}
		if updated.Name() != "Jane" {
			t.Errorf("ApplyAuthor() Name = %q, want %q", updated.Name(), "Jane")
		}
		if original.Name() != "John" {
			t.Errorf("ApplyAuthor() mutated original, Name = %q, want %q", original.Name(), "John")
		}
		if updated == original {
			t.Error("ApplyAuthor() returned same pointer as original, want new pointer")
		}
	})

	t.Run("empty changes return copy without error", func(t *testing.T) {
		original := entity.HydrateAuthor(5, "John")

		updated, err := original.ApplyAuthor(entity.AuthorChanges{})
		if err != nil {
			t.Fatalf("ApplyAuthor() unexpected error: %v", err)
		}
		if updated == nil {
			t.Fatal("ApplyAuthor() returned nil author")
		}
		if updated.Name() != "John" {
			t.Errorf("ApplyAuthor() Name = %q, want %q", updated.Name(), "John")
		}
		if updated == original {
			t.Error("ApplyAuthor() returned same pointer as original, want new pointer")
		}
	})
}

func TestApplyAuthor_Errors(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "empty name", arg: ""},
		{name: "spaces name", arg: "  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := entity.HydrateAuthor(5, "John")

			updated, err := original.ApplyAuthor(entity.AuthorChanges{Name: strPtr(tt.arg)})
			if err == nil {
				t.Fatal("ApplyAuthor() expected error, got nil")
			}
			if !errors.Is(err, entity.ErrValidationError) {
				t.Errorf("ApplyAuthor() error = %v, want ErrValidationError", err)
			}
			if updated != nil {
				t.Errorf("ApplyAuthor() author = %+v, want nil", updated)
			}
			if original.Name() != "John" {
				t.Errorf("ApplyAuthor() mutated original, Name = %q, want %q", original.Name(), "John")
			}
		})
	}
}
