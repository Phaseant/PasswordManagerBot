package repository

type Repository interface {
	Add(service, username, password, userID string) error
	Delete(serive, userID string) error
	Get(service, userID string) (string, string, error)
}
