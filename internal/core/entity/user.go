package entity

type User struct {
	ID    int64
	Name  string
	Email string
	Password string
	PasswordHash string
}

func (u *User) IsValid() bool {
	return u.Name != "" && u.Email != "" && u.Password != ""
}