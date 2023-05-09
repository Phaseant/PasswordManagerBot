package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repo interface {
	Add(service, username, password, userID string) error
	Delete(serive, username, userID string) error
	Get(service, userID string) ([]ServiceCreds, error)
}

type Repository struct {
	Repo
}

func New(db *mongo.Client, secretKey string) *Repository {
	return &Repository{
		Repo: DataRepository{
			db:        db,
			secretKey: secretKey,
		},
	}
}
