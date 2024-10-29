package internal

type ClientProfile struct {
	Email string
	Id string
	Name string
	Token string
}

// db mockup 
var Database = map[string]ClientProfile {
	"user1": {
		Email: "luis@luis.de",
		Id: "user1",
		Name: "Luis Thieme",
		Token: "123",
	},
	"user2": {
		Email: "max@max.de",
		Id: "user2",
		Name: "Max Zoladz",
		Token: "abc",
	},
}
