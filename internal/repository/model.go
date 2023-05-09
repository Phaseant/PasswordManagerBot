package repository

type Document struct {
	Service  string `bson:"service"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	UserID   string `bson:"userID"`
}
