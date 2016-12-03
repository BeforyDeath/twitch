package irc

type user struct {
	Name string
	Nick string
}

type capability struct {
	Broadcaster string
	Moderator   string
	Subscriber  string
	Turbo       string
	Premium     string
}
