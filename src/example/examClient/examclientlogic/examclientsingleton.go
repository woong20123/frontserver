package examclientlogic

//ExamSvrSingleton is
type ExamClientSingleton struct {
	objmanager *Objmanager
}

var instance *ExamClientSingleton = nil

// GetInstance is
func GetInstance() *ExamClientSingleton {
	if instance == nil {
		instance = newExamClientSingleton()
	}
	return instance
}

// newExamSvrSingleton is
func newExamClientSingleton() *ExamClientSingleton {
	ess := new(ExamClientSingleton)
	ess.objmanager = new(Objmanager)
	ess.objmanager.Intialize()

	return ess
}

// GetObjMgr is
func (s *ExamClientSingleton) GetObjMgr() *Objmanager {
	return s.objmanager
}
