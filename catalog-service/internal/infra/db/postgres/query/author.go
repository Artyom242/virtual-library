package query

const (
	CreateAuthor  = "INSERT INTO authors(name) VALUES($1) returning id, name"
	GetAuthorByID = "SELECT authors.id, authors.name FROM authors WHERE id = $1"
	GetAllAuthors = `
		SELECT authors.id, authors.name 
		FROM authors
		ORDER BY authors.id
		LIMIT $1
		OFFSET $2
		`
	GetOrCreateAuthor = `
		INSERT INTO authors(name) 
		VALUES($1)  
		ON CONFLICT (name) 
		    DO NOTHING 
			RETURNING id, name
	`
	DeleteAuthor = "DELETE FROM authors WHERE id = $1"
	UpdateAuthor = `UPDATE authors 
		SET name = COALESCE($1, name) 
		WHERE id = $2 
		RETURNING id, name`

	GetTotalAuthors       = `SELECT COUNT(*) FROM authors`
	GetAuthorsByName      = `SELECT id, name FROM authors WHERE name = $1`
	ExistsBooksByAuthorID = `SELECT EXISTS(SELECT 1 FROM books WHERE author_id = $1)`
)
