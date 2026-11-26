package query

const (
	CreateBook = `INSERT INTO books(title, year, author_id, description) 
		VALUES($1, $2, $3, $4) 
		RETURNING id, title, year, author_id, description
	`
	GetBookById = `SELECT id, title, year, author_id, description FROM books WHERE id = $1`
	GetAllBooks = `
		SELECT id, title, year, author_id, description
		FROM books
		ORDER BY id
		LIMIT $1
		OFFSET $2
	`
	DeleteBook = `DELETE FROM books WHERE id = $1`
	UpdateBook = `
		UPDATE books
		SET title     = COALESCE($1, title),
		    year      = COALESCE($2, year),
		    author_id = COALESCE($3, author_id),
		    description = COALESCE($4, description)
		WHERE id = $5
		RETURNING id, title, author_id, year, description
	`

	GetTotalBooks = `SELECT COUNT(*) FROM books`

	// Book_categories
	DeleteBookCategories = `DELETE FROM book_categories WHERE book_id = $1`
	SetCategories        = `
		INSERT INTO book_categories(book_id, category_id)
		VALUES($1, $2)
		RETURNING category_id
`
	GetCategoriesByBookID  = `SELECT category_id FROM book_categories WHERE book_id = $1`
	GetCategoriesByBookIDs = `SELECT book_id, category_id FROM book_categories WHERE book_id = ANY($1)`
)
