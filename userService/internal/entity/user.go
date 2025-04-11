package entity

type User struct {
	ID           string `json:"id" bson:"id"`
	Username     string `bson:"username" json:"username"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
}

type UserProfile struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
