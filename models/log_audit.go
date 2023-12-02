package models

import "time"

type LogAudit struct {
	Id           int64         `gorm:"column:id;type:bigint;not null;index:idx_pk_log_audit,unique;autoIncrement"`
	ApplId       string        `gorm:"column:applid;type:varchar(2);not null"`
	BranchId     string        `gorm:"column:branchid;type:varchar(4);not null"`
	DbName       string        `gorm:"column:dbname;type:varchar(50);not null"`
	TableNm      string        `gorm:"column:tablenm;type:varchar(50);not null"`
	TableKey     string        `gorm:"column:tablekey;type:varchar(255);not null"`
	ExecType     string        `gorm:"column:exectype;type:char;not null"`
	ExecUser     string        `gorm:"column:execuser;type:varchar(20);not null"`
	ExecDt       time.Time     `gorm:"column:execdt;type:date;not null"`
	IpAddr       string        `gorm:"column:ipaddr;type:varchar(50);not null"`
	SpvId        string        `gorm:"column:spvid;type:varchar(20)"`
	TimeStamp    time.Time     `gorm:"column:time_stamp;type:timestamp;not null"`
	LogAuditDtls []LogAuditDtl `gorm:"foreignKey:logid;references:id"`
}

func (l *LogAudit) TableName() string {
	return "log_audit"
}

type LogAuditDtl struct {
	LogId   int64  `gorm:"column:logid;type:bigint;not null;index:idx_pk_log_audit_dtl,unique"`
	Id      int64  `gorm:"column:id;type:bigint;not null;index:idx_pk_log_audit_dtl,unique;autoIncrement"`
	FieldNm string `gorm:"column:field;type:varchar(100);not null"`
	NewData string `gorm:"column:newdata;type:varchar(255);not null"`
	OldData string `gorm:"column:olddata;type:varchar(255)"`
}

func (l *LogAuditDtl) TableName() string {
	return "log_audit_dtl"
}
