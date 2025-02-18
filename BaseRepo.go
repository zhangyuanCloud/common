package common

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/sirupsen/logrus"
)

type BaseRepo struct {
	TableName string
	Log       *logrus.Entry
}

func (r *BaseRepo) InsertOne(o orm.Ormer, m interface{}) (int64, error) {
	if o == nil {
		o = orm.NewOrm()
	}
	insert, err := o.Insert(m)
	if err != nil {
		return 0, NewMsgError(CommonDbInsertError, err.Error())
	}
	return insert, nil
}

func (r *BaseRepo) InsertBatch(o orm.Ormer, bulk int, m interface{}) (int64, error) {
	if o == nil {
		o = orm.NewOrm()
	}
	return o.InsertMulti(bulk, m)
}

func (r *BaseRepo) ReadOne(m interface{}, cols ...string) error {
	return orm.NewOrm().Read(m, cols...)
}

func (r *BaseRepo) Update(o orm.Ormer, m interface{}, cols ...string) error {
	if o == nil {
		o = orm.NewOrm()
	}
	_, err := o.Update(m, cols...)
	if err != nil {
		return NewMsgError(CommonDbUpdateError, err.Error())
	}
	return nil
}

func (r *BaseRepo) UpdateByCondition(o orm.Ormer, cond *orm.Condition, param orm.Params) (int64, error) {
	if o == nil {
		o = orm.NewOrm()
	}
	if len(param) <= 0 {
		return 0, orm.ErrArgs
	}
	query := o.QueryTable(r.TableName)
	if cond != nil && !cond.IsEmpty() {
		query = query.SetCond(cond)
	}
	update, err := query.Update(param)
	if err != nil {
		return 0, NewMsgError(CommonDbUpdateError, err.Error())
	}
	return update, nil
}

func (r *BaseRepo) Delete(o orm.Ormer, m interface{}, cols ...string) error {
	if o == nil {
		o = orm.NewOrm()
	}
	_, err := o.Delete(m, cols...)
	return err
}

func (r *BaseRepo) DeleteByCondition(o orm.Ormer, cond *orm.Condition) error {
	if cond.IsEmpty() {
		return orm.ErrArgs
	}
	if o == nil {
		o = orm.NewOrm()
	}
	_, err := o.QueryTable(r.TableName).SetCond(cond).Delete()
	return err
}

func (r *BaseRepo) Count(cond *orm.Condition) int64 {
	query := orm.NewOrm().QueryTable(r.TableName).SetCond(cond)
	total, err := query.Count()
	if err != nil {
		return 0
	}
	return total
}

func (r *BaseRepo) List(cond *orm.Condition, sort string, container interface{}) (int64, error) {
	query := orm.NewOrm().QueryTable(r.TableName)
	if cond != nil {
		query = query.SetCond(cond)
	}
	total, err := query.Count()
	if err != nil {
		return 0, err
	}
	if len(sort) > 0 {
		query = query.OrderBy(sort)
	}
	_, err = query.All(container)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *BaseRepo) PageList(cond *orm.Condition, pageParam *BaseQueryParam, sort string, container interface{}) (int64, error) {

	query := orm.NewOrm().QueryTable(r.TableName).SetCond(cond)
	total, err := query.Count()
	if err != nil {
		return 0, err
	}

	if len(sort) > 0 {
		query = query.OrderBy(sort)
	}

	if pageParam.IsValid() {
		limit, offset := pageParam.GetLimit()
		query = query.Limit(limit).Offset(offset)
	}

	_, err = query.All(container)
	if err != nil {
		return 0, err
	}

	return total, nil
}
