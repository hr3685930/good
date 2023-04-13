package gorm

import (
	"gorm.io/gorm"
)

func (op OpentracingPlugin)beforeCreate(db *gorm.DB){
	op.injectBefor(db, _createOp)
}

func (op OpentracingPlugin)beforeUpdate(db *gorm.DB){
	op.injectBefor(db, _updateOp)
}

func (op OpentracingPlugin)beforeQuery(db *gorm.DB){
	op.injectBefor(db, _queryOp)
}

func (op OpentracingPlugin)beforeDelete(db *gorm.DB){
	op.injectBefor(db, _deleteOp)
}

func (op OpentracingPlugin)beforeRow(db *gorm.DB){
	op.injectBefor(db, _rowOp)
}

func (op OpentracingPlugin)beforeRaw(db *gorm.DB){
	op.injectBefor(db, _rawOp)
}

func (op OpentracingPlugin)after(db *gorm.DB){
	op.extractAfter(db)
}