package commons

type User struct {
	Name string
	Addr string
}

type Message struct {
	Content string
	Sender User
}

