package model

//struct for creating account
type User struct {
	UserName string `json:"userName" db:"user_name" validate:"userName"`
	Password string `json:"password" db:"password" validate:"password"`
	Prefect  string `json:"prefect" db:"prefect" validate:"prefect"`
	Gender   string `json:"gender" db:"gender" validate:"gender"`
}

//struct for login
type UserSimple struct {
	UserName string `json:"userName" db:"user_name" validate:"userName"`
	Password string `json:"password" db:"password" validate:"password"`
}

//struct for editing user info 
type UserUpdate struct {
	UserName    string `json:"userName" db:"user_name" validate:"userName"`
	Password    string `json:"password" validate:"password"`
	NewPassword string `json:"newPassword" validate:"password"`
	NeWPrefect  string `json:"newPrefect" validate:"prefect"`
	NeWGender   string `json:"newGender" validate:"gender"`
}

//struct for return value 
type UserWithoutPass struct {
	UserId   string `json:"userId" db:"user_id"`
	UserName string `json:"userName" db:"user_name" validate:"userName"`
	Prefect  string `json:"prefect" db:"prefect" validate:"prefect"`
	Gender   string `json:"gender" db:"gender" validate:"gender"`
}

//struct for return users according to user term
type UserList struct {
	HasNext bool           `json:"hasNext"`
	Users   *[]*UserWithoutPass `json:"users"`
}

//struct for delete info
type UserDelete struct {
	UserId	 string `json:"userId" db:"user_id"`
	Flag     int    `json:"flag" db:"flag"`
}