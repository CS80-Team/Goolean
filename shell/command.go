package shell

type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     func([]string) Status
}

type EarlyCommand struct {
	Name        string
	Description string
	Usage       string
	Handler     func()
}
