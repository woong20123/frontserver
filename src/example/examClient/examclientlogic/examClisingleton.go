package examclientlogic

//ExamCliSingleton is
type ExamCliSingleton struct {
	objmanager *Objmanager
}

var instance *ExamCliSingleton = nil

// GetInstance is
func GetInstance() *ExamCliSingleton {
	if instance == nil {
		instance = newExamCliSingleton()
	}
	return instance
}

// newExamSvrSingleton is
func newExamCliSingleton() *ExamCliSingleton {
	ess := new(ExamCliSingleton)
	ess.objmanager = new(Objmanager)
	ess.objmanager.Intialize()

	return ess
}

// GetObjMgr is
func (s *ExamCliSingleton) GetObjMgr() *Objmanager {
	return s.objmanager
}