package query

const (
	CreateCategory   = `INSERT INTO categories(name) VALUES($1) returning id, name`
	GetCategoryByID  = `SELECT categories.id, categories.name FROM categories WHERE id = $1`
	GetAllCategories = `
		SELECT categories.id, categories.name 
		FROM categories
		ORDER BY categories.id
		LIMIT $1
		OFFSET $2
		`
	DeleteCategory = `DELETE FROM categories WHERE id = $1`
	UpdateCategory = `UPDATE categories 
		SET name = COALESCE($1, name) 
		WHERE id = $2 
		RETURNING id, name
		`

	GetTotalCategories = `SELECT COUNT(*) FROM categories`
)
