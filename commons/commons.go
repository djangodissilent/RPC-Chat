package commons

type Client struct {
	Name string
	Addr string
}

type Message struct {
	Content string
	Sender Client
}

