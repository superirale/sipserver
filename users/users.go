package users

// User structure
type User struct {
	Username string
	Password string
}

// MakeUsers a factory Adding default users
func MakeUsers() map[string]*User {
	users := make(map[string]*User)
	users["superirale"] = newUser("superirale", "omokhudu")
	users["usmanirale"] = newUser("usmanirale", "omokhudu1987")
	users["hackangel"] = newUser("hackangel", "omokhudu")

	return users
}

//newUser a constructor for creating user accounts
func newUser(username string, pswd string) *User {
	user := new(User)
	user.Username = username
	user.Password = pswd

	return user
}