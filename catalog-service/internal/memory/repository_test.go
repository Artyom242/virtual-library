package memory

import (
	"reflect"
	"testing"
	"virtual-library/catalog-service/internal/catalog"
)

func TestMemoryBookRepository_Create(t *testing.T) {
	tests := []struct {
		name   string
		data   catalog.Book
		bookID int
	}{
		{
			name: "Create new first book",
			data: catalog.Book{
				Author: "Tom Wednesday",
				Title:  "One horse",
				Year:   2001,
			},
			bookID: 1,
		},
		{
			name: "Create new second book",
			data: catalog.Book{
				Author: "John Doe",
				Title:  "The Big Buck",
				Year:   2006,
			},
			bookID: 2,
		},
	}

	repo := NewMemoryBookRepository()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Create(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != test.bookID {
				t.Errorf("got ID %v, want %v", got.ID, test.bookID)
			}
			if got.Author != test.data.Author {
				t.Errorf("got Author %v, want %v", got.Author, test.data.Author)
			}
			if got.Title != test.data.Title {
				t.Errorf("got Title %v, want %v", got.Title, test.data.Title)
			}
			if got.Year != test.data.Year {
				t.Errorf("got Year %v, want %v", got.Year, test.data.Year)
			}
		})
	}
}
func TestMemoryBookRepository_GetByID(t *testing.T) {
	repo := NewMemoryBookRepository()

	book, err := repo.Create(catalog.Book{
		Author: "John Doe",
		Title:  "The Big Buck",
		Year:   2006,
	})

	if err != nil {
		t.Fatalf("failed to create book %v", err)
	}

	tests := []struct {
		name    string
		want    *catalog.Book
		wantErr bool
		bookID  int
	}{
		{
			name:   "Get existing book",
			want:   &book,
			bookID: book.ID,
		},
		{
			name:    "Get not existing book",
			want:    nil,
			wantErr: true,
			bookID:  99,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.GetByID(test.bookID)

			if (err != nil) != test.wantErr {
				t.Fatalf("error getting book %v", err)
			}

			if got == nil && test.want != nil {
				t.Fatal("expected book, got nil")
			}

			if got != nil && test.want == nil {
				t.Fatalf("expected nil, got %v", got)
			}

			if got != nil && test.want != nil && *got != *test.want {
				t.Errorf("got %+v, want %+v", *got, *test.want)
			}
		})
	}
}
func TestMemoryBookRepository_GetAll(t *testing.T) {
	tests := []struct {
		name  string
		books []catalog.Book
		want  []catalog.Book
	}{
		{
			name: "Get not existing book",
			want: make([]catalog.Book, 0),
		},
		{
			name: "Get one book",
			books: []catalog.Book{
				{
					Author: "John Doe",
					Title:  "The Big Buck",
					Year:   2006,
				},
			},
			want: []catalog.Book{
				{
					ID:     1,
					Author: "John Doe",
					Title:  "The Big Buck",
					Year:   2006,
				},
			},
		},
		{
			name: "Get all books",
			books: []catalog.Book{
				{
					Author: "John Doe",
					Title:  "The Big Buck",
					Year:   2006,
				},
				{
					Author: "Tom Wednesday",
					Title:  "One horse",
					Year:   2005,
				},
			},
			want: []catalog.Book{
				{
					ID:     1,
					Author: "John Doe",
					Title:  "The Big Buck",
					Year:   2006,
				},
				{
					ID:     2,
					Author: "Tom Wednesday",
					Title:  "One horse",
					Year:   2005,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoTest := NewMemoryBookRepository()

			for i := 0; i < len(test.books); i++ {
				_, err := repoTest.Create(test.books[i])
				if err != nil {
					t.Fatal(err)
				}
			}

			got, err := repoTest.GetAll()

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %+v, want %+v", got, test.want)
			}
		})
	}
}
func TestMemoryBookRepository_Delete(t *testing.T) {
	tests := []struct {
		name         string
		isCreateBook bool
		wantErr      bool
	}{
		{
			name:         "Delete existing book",
			isCreateBook: true,
		},
		{
			name:    "Delete not existing book",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewMemoryBookRepository()
			bookID := 99
			if test.isCreateBook {
				book, err := repo.Create(catalog.Book{
					Author: "John Doe",
					Title:  "The Big Buck",
					Year:   2006,
				})
				if err != nil {
					t.Fatalf("failed to create book %v", err)
				}

				bookID = book.ID
			}

			err := repo.Delete(bookID)
			if (err != nil) != test.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, test.wantErr)
			}

			got, err := repo.GetByID(bookID)
			if err == nil {
				t.Fatal("expected error when getting deleted or non-existing book")
			}
			if got != nil {
				t.Fatalf("got %+v, want nil", got)
			}
		})
	}
}
