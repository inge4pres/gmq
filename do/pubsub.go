package gmq

type PubSub struct {
	Prio  int
	Name  string
	Queue QManager
}
