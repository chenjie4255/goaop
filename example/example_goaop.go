// This file is generated by ifmeasure, DO NOT EDIT IT.
// see (github.com/chenjie4255/goaop) 
package example
import (
	"github.com/chenjie4255/goaop"
	
ppp "github.com/chenjie4255/goaop/example/param"
)

type measureUserDB struct {
	UserDB
	builder goaop.PointcutBuilder
}

func NewMeasureUserDB(o UserDB, builder goaop.PointcutBuilder) UserDB {
	return &measureUserDB{o, builder}
}


func (m *measureUserDB) GetUserCount() (int, error) {
	pointcut := m.builder.Build("GetUserCount")
	if pointcut == nil {
		return m.GetUserCount()
	} else {
		pointcut.OnEntry()
		 r0,err := m.GetUserCount()
		pointcut.OnReturn(err)
		return r0,err 
	}
}

func (m *measureUserDB) SetUserScore(userID string, score ppp.Score) {
	pointcut := m.builder.Build("SetUserScore")
	if pointcut == nil {
		m.SetUserScore(userID ,score)
	} else {
		pointcut.OnEntry()
		m.SetUserScore(userID ,score)
		pointcut.OnReturn(nil)
		
	}
}

func (m *measureUserDB) SetUserScores(userID string, scores []ppp.Score) error {
	pointcut := m.builder.Build("SetUserScores")
	if pointcut == nil {
		return m.SetUserScores(userID ,scores)
	} else {
		pointcut.OnEntry()
		 err := m.SetUserScores(userID ,scores)
		pointcut.OnReturn(err)
		return err 
	}
}

func (m *measureUserDB) RemoveUser(userIDs ...string) error {
	pointcut := m.builder.Build("RemoveUser")
	if pointcut == nil {
		return m.RemoveUser(userIDs...)
	} else {
		pointcut.OnEntry()
		 err := m.RemoveUser(userIDs...)
		pointcut.OnReturn(err)
		return err 
	}
}

func (m *measureUserDB) UpdateUserBatch(userIDs []string, scores []ppp.Score) (int, int) {
	pointcut := m.builder.Build("UpdateUserBatch")
	if pointcut == nil {
		return m.UpdateUserBatch(userIDs ,scores)
	} else {
		pointcut.OnEntry()
		 r0,r1 := m.UpdateUserBatch(userIDs ,scores)
		
		pointcut.OnReturn(nil)
		
		return r0,r1 
	}
}



