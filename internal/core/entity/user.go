package entity

type User struct {
	ID    int64
	Name  string
	Email string
}

func (u *User) IsValid() bool {
	return u.Name != "" && u.Email != ""
}