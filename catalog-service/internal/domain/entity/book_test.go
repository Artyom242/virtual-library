package entity_test

import (
	"errors"
	"reflect"
	"testing"

	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/domain/entity"
)

func strPtr(s string) *string { return &s }

func intPtr(v int) *int { return &v }

func intsPtr(v ...int) *[]int {
	if v == nil {
		v = []int{}
	}

	return &v
}

type bookState struct {
	id          int
	title       string
	year        int
	authorID    int
	description *string
	categories  []int
}

func snapshot(b *entity.Book) bookState {
	return bookState{
		id:          b.ID(),
		title:       b.Title(),
		year:        b.Year(),
		authorID:    b.AuthorID(),
		description: b.Description(),
		categories:  b.Categories(),
	}
}

func mustNewBook(t *testing.T, title string, year, authorID int, description *string, categories *[]int) *entity.Book {
	t.Helper()

	book, err := entity.NewBook(title, year, authorID, description, categories)
	if err != nil {
		t.Fatalf("NewBook() unexpected error: %v", err)
	}

	return book
}

func TestNewBook_Success(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		year        int
		authorID    int
		description *string
		categories  *[]int
		want        bookState
	}{
		{
			name:     "nil categories become empty slice",
			title:    "Go",
			year:     2020,
			authorID: 2,
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{},
			},
		},
		{
			name:        "full data",
			title:       "Go",
			year:        2020,
			authorID:    2,
			description: strPtr("ok"),
			categories:  intsPtr(1, 2),
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: strPtr("ok"),
				categories:  []int{1, 2},
			},
		},
		{
			name:     "lower bounds are valid",
			title:    "Go",
			year:     1,
			authorID: 1,
			want: bookState{
				id:          0,
				title:       "Go",
				year:        1,
				authorID:    1,
				description: nil,
				categories:  []int{},
			},
		},
		{
			name:       "empty categories slice stays empty",
			title:      "Go",
			year:       2020,
			authorID:   2,
			categories: intsPtr(),
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := entity.NewBook(tt.title, tt.year, tt.authorID, tt.description, tt.categories)
			if err != nil {
				t.Fatalf("NewBook() unexpected error: %v", err)
			}
			if book == nil {
				t.Fatal("NewBook() returned nil book")
			}

			got := snapshot(book)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBook() state = %+v, want %+v", got, tt.want)
			}
			if book.Categories() == nil {
				t.Error("Categories() returned nil slice, want non-nil")
			}
		})
	}
}

func TestNewBook_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		year        int
		authorID    int
		description *string
		categories  *[]int
	}{
		{name: "empty title", title: "", year: 2020, authorID: 2},
		{name: "spaces title", title: "   ", year: 2020, authorID: 2},
		{name: "whitespace title", title: "\t\n", year: 2020, authorID: 2},
		{name: "zero year", title: "Go", year: 0, authorID: 2},
		{name: "negative year", title: "Go", year: -1, authorID: 2},
		{name: "zero author id", title: "Go", year: 2020, authorID: 0},
		{name: "negative author id", title: "Go", year: 2020, authorID: -1},
		{name: "spaces description", title: "Go", year: 2020, authorID: 2, description: strPtr("  ")},
		{name: "whitespace description", title: "Go", year: 2020, authorID: 2, description: strPtr("\t")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := entity.NewBook(tt.title, tt.year, tt.authorID, tt.description, tt.categories)
			if err == nil {
				t.Fatal("NewBook() expected error, got nil")
			}
			if book != nil {
				t.Errorf("NewBook() book = %+v, want nil", book)
			}
			if !errors.Is(err, entity.ErrValidationError) {
				t.Errorf("NewBook() error = %v, want ErrValidationError", err)
			}
			if errors.Is(err, entity.ErrCategoryNotUnique) {
				t.Errorf("NewBook() error = %v, must not be ErrCategoryNotUnique", err)
			}
		})
	}
}

func TestNewBook_CategoryNotUnique(t *testing.T) {
	tests := []struct {
		name       string
		categories *[]int
	}{
		{name: "duplicate pair", categories: intsPtr(1, 1)},
		{name: "duplicate among many", categories: intsPtr(1, 2, 3, 2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := entity.NewBook("Go", 2020, 2, nil, tt.categories)
			if err == nil {
				t.Fatal("NewBook() expected error, got nil")
			}
			if book != nil {
				t.Errorf("NewBook() book = %+v, want nil", book)
			}
			if !errors.Is(err, entity.ErrCategoryNotUnique) {
				t.Errorf("NewBook() error = %v, want ErrCategoryNotUnique", err)
			}
			if errors.Is(err, entity.ErrValidationError) {
				t.Errorf("NewBook() error = %v, must not wrap ErrValidationError", err)
			}
		})
	}
}

func TestHydrateBook(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		title       string
		year        int
		authorID    int
		description *string
		categories  *[]int
		want        bookState
	}{
		{
			name:     "nil categories become empty slice",
			id:       5,
			title:    "Go",
			year:     2020,
			authorID: 2,
			want: bookState{
				id:          5,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{},
			},
		},
		{
			name:        "invalid data is kept without validation",
			id:          7,
			title:       "",
			year:        0,
			authorID:    0,
			description: strPtr("x"),
			categories:  intsPtr(3),
			want: bookState{
				id:          7,
				title:       "",
				year:        0,
				authorID:    0,
				description: strPtr("x"),
				categories:  []int{3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := entity.HydrateBook(tt.id, tt.title, tt.year, tt.authorID, tt.description, tt.categories)
			if book == nil {
				t.Fatal("HydrateBook() returned nil")
			}

			got := snapshot(book)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HydrateBook() state = %+v, want %+v", got, tt.want)
			}
			if book.Categories() == nil {
				t.Error("Categories() returned nil slice, want non-nil")
			}
		})
	}
}

func TestApplyChanges_Success(t *testing.T) {
	baseCategories := intsPtr(1, 2)

	tests := []struct {
		name    string
		start   *entity.Book
		changes entity.BookChanges
		want    bookState
	}{
		{
			name:  "title only keeps categories",
			start: mustNewBook(t, "Go", 2020, 2, nil, baseCategories),
			changes: entity.BookChanges{
				Title: strPtr("New"),
			},
			want: bookState{
				id:          0,
				title:       "New",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{1, 2},
			},
		},
		{
			name:  "replace categories",
			start: mustNewBook(t, "Go", 2020, 2, nil, baseCategories),
			changes: entity.BookChanges{
				CategoriesIDs: intsPtr(3),
			},
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{3},
			},
		},
		{
			name:  "clear categories with empty slice",
			start: mustNewBook(t, "Go", 2020, 2, nil, baseCategories),
			changes: entity.BookChanges{
				CategoriesIDs: intsPtr(),
			},
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: nil,
				categories:  []int{},
			},
		},
		{
			name:  "update year and author",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				Year:     intPtr(2021),
				AuthorID: intPtr(3),
			},
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2021,
				authorID:    3,
				description: nil,
				categories:  []int{},
			},
		},
		{
			name:  "set description",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				Description: strPtr("d"),
			},
			want: bookState{
				id:          0,
				title:       "Go",
				year:        2020,
				authorID:    2,
				description: strPtr("d"),
				categories:  []int{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.start.ApplyChanges(tt.changes); err != nil {
				t.Fatalf("ApplyChanges() unexpected error: %v", err)
			}

			got := snapshot(tt.start)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyChanges() state = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestApplyChanges_Errors(t *testing.T) {
	tests := []struct {
		name    string
		start   *entity.Book
		changes entity.BookChanges
		wantErr error
	}{
		{
			name:    "empty patch",
			start:   mustNewBook(t, "Go", 2020, 2, nil, intsPtr(1, 2)),
			changes: entity.BookChanges{},
			wantErr: app.ErrEmptyBookPatch,
		},
		{
			name:  "empty title",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				Title: strPtr(""),
			},
			wantErr: entity.ErrValidationError,
		},
		{
			name:  "zero year",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				Year: intPtr(0),
			},
			wantErr: entity.ErrValidationError,
		},
		{
			name:  "zero author id",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				AuthorID: intPtr(0),
			},
			wantErr: entity.ErrValidationError,
		},
		{
			name:  "whitespace description",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				Description: strPtr(" "),
			},
			wantErr: entity.ErrValidationError,
		},
		{
			name:  "duplicate categories",
			start: mustNewBook(t, "Go", 2020, 2, nil, nil),
			changes: entity.BookChanges{
				CategoriesIDs: intsPtr(1, 1),
			},
			wantErr: entity.ErrCategoryNotUnique,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := snapshot(tt.start)

			err := tt.start.ApplyChanges(tt.changes)
			if err == nil {
				t.Fatal("ApplyChanges() expected error, got nil")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApplyChanges() error = %v, want %v", err, tt.wantErr)
			}

			after := snapshot(tt.start)
			if !reflect.DeepEqual(after, before) {
				t.Errorf("ApplyChanges() mutated book on error: before = %+v, after = %+v", before, after)
			}
		})
	}
}
