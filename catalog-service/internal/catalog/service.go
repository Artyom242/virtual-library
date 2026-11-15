package catalog

type BookService struct {
	repo BookRepository
}

func NewService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}
func (s *BookService) CreateBook(book Book) (Book, error) {
	return s.repo.Create(book)
}

func (s *BookService) GetByID(id int) (*Book, error) {
	return s.repo.GetByID(id)
}

func (s *BookService) GetAll() ([]Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) Delete(bookID int) error {
	return s.repo.Delete(bookID)
}
