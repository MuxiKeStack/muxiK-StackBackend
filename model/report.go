package model

// the table name for report model, default will be reports
func (rm *ReportModel) TableName() string {
	return "report"
}

// create a new report
func (rm *ReportModel) Create() error {
	d := DB.Self.Create(rm)
	return d.Error
}

// update a report
func (rm *ReportModel) Update() error {
	d := DB.Self.Update(rm)
	return d.Error
}

// check the report is existed
func ReportExisted(eid uint64, uid uint32) bool {
	if DB.Self.Table("report").Where("evaluation_id = ? AND user_id = ?", eid, uid).RecordNotFound() {
		return false
	}
	return true
}

// get an evaluation be reported total
func CountEvaluationBeReportedTimes(eid uint64) int {
	var count int
	DB.Self.Table("report").Where("evaluation_id = ? AND pass = false", eid).Count(&count)
	return count
}

// get the all report of an evaluation
func GetAllReportOfEvaluation(eid uint64) ([]ReportModel, error) {
	var reports []ReportModel
	d := DB.Self.Table("report").Where("evaluation_id = ? AND pass = false", eid).Find(&reports)
	if d.Error != nil {
		return nil, d.Error
	}
	return reports, nil
}
