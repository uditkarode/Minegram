package modules

import (
	"regexp"
)

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
var feed = make(chan consoleEvent)
var logFeed = make(chan string)

func emit(k eventKind, str string) {
	if k == eventLog {
		logFeed <- str
	} else {
		for i := 0; i < listenerCount; i++ {
			feed <- consoleEvent{line: str, kind: k}
		}
	}
}

func listenLog(action func(string)) {
	for {
		action(<-logFeed)
	}
}

/* ALWAYS use this in a separate goroutine! */
func listen(k eventKind, action func(str string)) {
	listenerCount++
	for {
		event := <-feed
		if event.kind == k {
			action(event.line)
		}
	}
}

/* TEMP REMOVE LATER */
var chatRegex = regexp.MustCompile(`: <(.+)> (.+)`)
var joinRegex = regexp.MustCompile(`: (.+) joined the game`)
var joinRegexSpigotPaper = regexp.MustCompile(`: UUID of player (.+) is .*`)
var leaveRegex = regexp.MustCompile(`: (.+) left the game`)
var advancementRegex = regexp.MustCompile(`: (.+) has made the advancement (.+)`)

/* death regex taken from https://github.com/trgwii/TeMiCross/blob/master/client/parser/default/messages/death.js */
var deathRegex = regexp.MustCompile(`: (.+) (was (shot by .+|shot off (some vines|a ladder) by .+|pricked to death|stabbed to death|squished too much|blown up by .+|killed by .+|doomed to fall by .+|blown from a high place by .+|squashed by .+|burnt to a crisp whilst fighting .+|roasted in dragon breath( by .+)?|struck by lightning( whilst fighting .+)?|slain by .+|fireballed by .+|killed trying to hurt .+|impaled by .+|speared by .+|poked to death by a sweet berry bush( whilst trying to escape .+)?|pummeled by .+)|hugged a cactus|walked into a cactus whilst trying to escape .+|drowned( whilst trying to escape .+)?|suffocated in a wall( whilst fighting .+)?|experienced kinetic energy( whilst trying to escape .+)?|removed an elytra while flying( whilst trying to escape .+)?|blew up|hit the ground too hard( whilst trying to escape .+)?|went up in flames|burned to death|walked into fire whilst fighting .+|went off with a bang( whilst fighting .+)?|tried to swim in lava(( while trying)? to escape .+)?|discovered floor was lava|walked into danger zone due to .+|got finished off by .+|starved to death|didn't want to live in the same world as .+|withered away( whilst fighting .+)?|died( because of .+)?|fell (from a high place( and fell out of the world)?|off a ladder|off to death( whilst fighting .+)?|off some vines|out of the water|into a patch of fire|into a patch of cacti|too far and was finished by .+|out of the world))$`)

var timeRegex = regexp.MustCompile(`: The time is (.+)`)
var entityPosRegex = regexp.MustCompile(`: .+ has the following entity data: \[(.+)d, (.+)d, (.+)d\]`)
var simplifiedEPRegex = regexp.MustCompile(`: .+ has the following entity data: \[(.+)\..*d, (.+)\..*d, (.+)\..*d\]`)
var simpleOutputRegex = regexp.MustCompile(`.*: (.+)`)
var dimensionRegex = regexp.MustCompile(`.*has the following entity data: "(minecraft:.+)"`)
var gameTypeRegex = regexp.MustCompile(`.*has the following entity data: (.+)`)
var genericOutputRegex = regexp.MustCompile(`(\[.+\]) (\[.+\]): (.+)`)