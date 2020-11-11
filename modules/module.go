package modules

type eventKind string

var listenerCount = 0

const (
	eventLog         eventKind = "EVENT_LOG"
	eventChat        eventKind = "EVENT_CHAT"
	eventJoin        eventKind = "EVENT_JOIN"
	eventLeave       eventKind = "EVENT_LEAVE"
	eventAdvancement eventKind = "EVENT_ADVANCEMENT"
	eventDeath       eventKind = "EVENT_DEATH"
)

type consoleEvent struct {
	line string
	kind eventKind
}

var err error
var logFeed = make(chan string)
