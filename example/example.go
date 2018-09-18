package example

import (
	ppp "github.com/chenjie4255/goaop/example/param"
)

//go:generate goaop -f=$GOFILE

type intParam struct {
}

// UserDB user's db
// @ifmeasure
type UserDB interface {
	GetUserCount() (int, error)
	SetUserScore(userID string, score ppp.Score)
	SetUserScores(userID string, scores []ppp.Score) error
	RemoveUser(userIDs ...string) error
	UpdateUserBatch(userIDs []string, scores []ppp.Score) (int, int)
}
