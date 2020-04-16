package examsvrsingleton

import "example/examFrontServer/serverobjmanager"

type ExamSvrSingleton struct {
	objmanager *serverobjmanager.SvrObjMgr
}

var instance *ExamSvrSingleton = nil

// GetInstance is
func GetInstance() *ExamSvrSingleton {
	if instance == nil {
		instance = newExamSvrSingleton()
	}
	return instance
}

// newExamSvrSingleton is
func newExamSvrSingleton() *ExamSvrSingleton {
	ess := new(ExamSvrSingleton)
	ess.objmanager = new(serverobjmanager.SvrObjMgr)
	ess.objmanager.Initialize()

	return ess
}
