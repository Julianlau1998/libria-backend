package models

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	Nickname      string `json:"nickname"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

// type UserDB struct {
// 	ID       string
// 	Username sql.NullString
// 	Password sql.NullString
// }

// func (dbV *UserDB) GetUser() (u User) {
// 	u.ID = dbV.ID
// 	u.Username = utility.GetStringValue(dbV.Username)
// 	u.Password = utility.GetStringValue(dbV.Password)
// 	return u
// }
