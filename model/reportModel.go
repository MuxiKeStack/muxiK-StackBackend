package model

/*
CREATE TABLE `report` (
  `id`            INT unsigned NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT          NOT NULL,
  `user_id`       INT          NOT NULL,
  `pass`          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "举报审核是否通过",
  `reason`        VARCHAR(200) NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;
*/

type ReportModel struct {
	Id           uint64 `gorm:"column:id; primary_key"`
	EvaluationId uint64 `gorm:"column:evaluation_id"`
	UserId       uint64 `gorm:"column:user_id"`
	Pass         bool   `gorm:"column:pass"`
	Reason       string `gorm:"column:reason"`
}
