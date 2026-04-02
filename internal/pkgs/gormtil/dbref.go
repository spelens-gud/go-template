package gormtil

import (
	"context"
	"strings"

	"{{.ProjectName}}/internal/pkgs/gotil"

	"git.bestfulfill.tech/devops/go-core/implements/gormx"
	"git.bestfulfill.tech/devops/go-core/interfaces/isql"
	"git.bestfulfill.tech/devops/go-core/kits/kerror/errorx"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// RefPlucker 引用关系
type RefPlucker interface {
	// A2ManyB 通过 A 获取所有 B
	A2ManyB(ctx context.Context, db *gorm.DB, a ...int) (b []int, err error)
	// B2ManyA 通过 B 获取所有 A
	B2ManyA(ctx context.Context, db *gorm.DB, b ...int) (a []int, err error)
	// ExistAnB 通过 A 和 B 判断是否存在
	ExistAnB(db *gorm.DB, a, b int) (exist bool, err error)
	// ExistManyAnOneB 通过 A 和 B 批量判断是否存在
	ExistManyAnOneB(db *gorm.DB, a []int, b int) (exist bool, err error)
	// ExistManyBnOneA 通过 B 和 A 批量判断是否存在
	ExistManyBnOneA(db *gorm.DB, b []int, a int) (exist bool, err error)
	// MappingSliceByA 通过 A 批量获取所有 B
	MappingSliceByA(ctx context.Context, db *gorm.DB, a []int) (ret map[int][]int, err error)
	// MappingSliceByB 通过 B 批量获取所有 A
	MappingSliceByB(ctx context.Context, db *gorm.DB, b []int) (ret map[int][]int, err error)
	// InsertRef 插入引用关系
	InsertRef(ctx context.Context, db *gorm.DB, a, b int) (err error)
	// UpdateRefByA 通过A更新引用关系
	UpdateRefByA(ctx context.Context, db *gorm.DB, a int, b []int) (err error)
	// UpdateRefByB 通过B更新引用关系
	UpdateRefByB(ctx context.Context, db *gorm.DB, b int, a []int) (err error)
	// CountByAConfig 获取 A 的引用关系数量
	CountByAConfig() DBMappingLoaderConfig
	// CountByBConfig 获取 B 的引用关系数量
	CountByBConfig() DBMappingLoaderConfig
	// DeleteRef 删除引用关系
	DeleteRef(ctx context.Context, db *gorm.DB, a, b int) (err error)
	// HasAnyData 判断是否有数据
	HasAnyData(db *gorm.DB) (has bool, err error)
	// Dump 获取所有数据
	Dump(ctx context.Context, db *gorm.DB) (data []RefData, err error)
	// Import 导入数据
	Import(ctx context.Context, db *gorm.DB, data []RefData) (err error)
}

// IntFromRow 从行中解析 int (如 id)
func IntFromRow(row any) (int, error) {
	return cast.ToIntE(row)
}

// RefData 引用关系数据
type RefData struct {
	A int // 引用关系 A
	B int // 引用关系 B
}

// refPlucker 实现 RefPlucker 接口
type refPlucker struct {
	FieldA       string            // 引用关系 A 的字段名
	FieldB       string            // 引用关系 B 的字段名
	Table        string            // 引用关系表名
	ExtKV        map[string]string // 额外的字段
	extCondition string            // 额外的条件
}

// NewRefPluck 创建 RefPlucker
func NewRefPluck(table, fieldA, fieldB string, extKV map[string]string) RefPlucker {
	return &refPlucker{
		FieldA: fieldA,
		FieldB: fieldB,
		Table:  table,
		ExtKV:  extKV,
	}
}

// A2ManyB 通过 A 获取所有 B
func (r *refPlucker) A2ManyB(ctx context.Context, db *gorm.DB, a ...int) (b []int, err error) {
	if len(a) == 0 {
		err = errorx.ErrInvalidArgument("无效的查询参数")
		return
	}

	loader := r.refSliceLoaderByA(db)
	m, err := loader.LoadMany(ctx, a)
	if err != nil {
		return nil, err
	}
	for _, s := range m {
		b = append(b, s...)
	}
	return gotil.NewSetWithSlice(b).Slice(), nil
}

// B2ManyA 通过 B 获取所有 A
func (r *refPlucker) B2ManyA(ctx context.Context, db *gorm.DB, b ...int) (a []int, err error) {
	if len(b) == 0 {
		err = errorx.ErrInvalidArgument("无效的查询参数")
		return
	}

	loader := r.refSliceLoaderByB(db)
	m, err := loader.LoadMany(ctx, b)
	if err != nil {
		return nil, err
	}
	for _, s := range m {
		a = append(a, s...)
	}
	return gotil.NewSetWithSlice(a).Slice(), nil
}

// ExistAnB 通过 A 和 B 判断是否存在
func (r *refPlucker) ExistAnB(db *gorm.DB, a, b int) (exist bool, err error) {
	return r.existWithConditions(db,
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldA+" = ?", a) },
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldB+" = ?", b) },
	)
}

// ExistManyAnOneB 通过 A 和 B 批量判断是否存在
func (r *refPlucker) ExistManyAnOneB(db *gorm.DB, a []int, b int) (exist bool, err error) {
	if len(a) == 0 {
		return
	}
	return r.existWithConditions(db,
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldA+" IN ?", a) },
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldB+" = ?", b) },
	)
}

// ExistManyBnOneA 通过 B 和 A 批量判断是否存在
func (r *refPlucker) ExistManyBnOneA(db *gorm.DB, b []int, a int) (exist bool, err error) {
	if len(b) == 0 {
		return
	}
	return r.existWithConditions(db,
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldB+" IN ?", b) },
		func(q *gorm.DB) *gorm.DB { return q.Where(r.FieldA+" = ?", a) },
	)
}

// MappingSliceByA 通过 A 批量获取所有 B
func (r *refPlucker) MappingSliceByA(ctx context.Context, db *gorm.DB, a []int) (ret map[int][]int, err error) {
	if len(a) == 0 {
		return nil, nil
	}
	loader := r.refSliceLoaderByA(db)
	return loader.LoadMany(ctx, a)
}

// MappingSliceByB 通过 B 批量获取所有 A
func (r *refPlucker) MappingSliceByB(ctx context.Context, db *gorm.DB, b []int) (ret map[int][]int, err error) {
	if len(b) == 0 {
		return nil, nil
	}
	loader := r.refSliceLoaderByB(db)
	return loader.LoadMany(ctx, b)
}

// InsertRef 插入引用关系
func (r *refPlucker) InsertRef(ctx context.Context, db *gorm.DB, a, b int) (err error) {
	keyFields := []string{r.FieldA, r.FieldB}
	extFields, extValue := r.extKVFields()

	if err = gormx.NewWithGormDB(db).BatchInsert(ctx, r.Table, []struct{}{{}}, func(i int, insert func(...any)) {
		insert(append([]any{a, b}, extValue...)...)
	}, append(keyFields, extFields...), isql.IsIgnore()); err != nil {
		return
	}
	return
}

// UpdateRefByA 通过A更新引用关系
func (r *refPlucker) UpdateRefByA(ctx context.Context, db *gorm.DB, a int, b []int) (err error) {
	return r.updateRef(ctx, db, a, b, r.FieldA, r.FieldB)
}

// UpdateRefByB 通过B更新引用关系
func (r *refPlucker) UpdateRefByB(ctx context.Context, db *gorm.DB, b int, a []int) (err error) {
	return r.updateRef(ctx, db, b, a, r.FieldB, r.FieldA)
}

// CountByAConfig 获取 A 的引用关系数量
func (r *refPlucker) CountByAConfig() DBMappingLoaderConfig {
	return DBMappingLoaderConfig{
		TableName:   r.Table,
		KeyColumn:   r.FieldA,
		ValueColumn: "count(*)",
		Where:       r.getExtConditionSlice(),
	}
}

// CountByBConfig 获取 B 的引用关系数量
func (r *refPlucker) CountByBConfig() DBMappingLoaderConfig {
	return DBMappingLoaderConfig{
		TableName:   r.Table,
		KeyColumn:   r.FieldB,
		ValueColumn: "count(*)",
		Where:       r.getExtConditionSlice(),
	}
}

// DeleteRef 删除引用关系
func (r *refPlucker) DeleteRef(ctx context.Context, db *gorm.DB, a, b int) (err error) {
	err = r.table(db).Where(r.FieldA+" = ?", a).Where(r.FieldB+" = ?", b).Delete(nil).Error
	return
}

// HasAnyData 判断是否有数据
func (r *refPlucker) HasAnyData(db *gorm.DB) (has bool, err error) {
	var ret []RefData
	err = r.table(db).Select(r.FieldA+" AS a", r.FieldB+" AS b").Limit(1).Find(&ret).Error
	if err != nil {
		return
	}
	has = len(ret) > 0
	return
}

// Dump 获取所有数据
func (r *refPlucker) Dump(ctx context.Context, db *gorm.DB) (data []RefData, err error) {
	err = r.table(db).Select(r.FieldA+" AS a", r.FieldB+" AS b").Find(&data).Error
	return
}

// Import 导入数据
func (r *refPlucker) Import(ctx context.Context, db *gorm.DB, data []RefData) (err error) {
	keyFields := []string{r.FieldA, r.FieldB}
	extFields, extValue := r.extKVFields()

	if err = gormx.NewWithGormDB(db).BatchInsert(ctx, r.Table, data, func(i int, insert func(...interface{})) {
		insert(append([]interface{}{data[i].A, data[i].B}, extValue...)...)
	}, append(keyFields, extFields...), isql.IsIgnore()); err != nil {
		return
	}
	return
}

// applyExtConditions 将 ExtKV 应用到 GORM 查询, 用于 Exist 等需要额外条件的场景
func (r *refPlucker) applyExtConditions(query *gorm.DB) *gorm.DB {
	for k, v := range r.ExtKV {
		query = query.Where(k+" = ?", v)
	}
	return query
}

// existWithConditions 使用 GORM 语法检查是否存在记录
func (r *refPlucker) existWithConditions(db *gorm.DB, conditions ...func(*gorm.DB) *gorm.DB) (exist bool, err error) {
	query := db.Table(r.Table)
	query = r.applyExtConditions(query)

	for _, cond := range conditions {
		query = cond(query)
	}
	var count int64
	err = query.Limit(1).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// refSliceLoaderByA 批量获取 A 的引用关系
func (r *refPlucker) refSliceLoaderByA(db *gorm.DB) *DBBatchSliceLoader[int, int] {
	cfg := DBMappingLoaderConfig{
		TableName:   r.Table,
		KeyColumn:   r.FieldA,
		ValueColumn: r.FieldB,
		Where:       r.getExtConditionSlice(),
	}
	return NewDBBatchSliceLoader[int, int](db, cfg, IntFromRow, IntFromRow)
}

// refSliceLoaderByB 批量获取 B 的引用关系
func (r *refPlucker) refSliceLoaderByB(db *gorm.DB) *DBBatchSliceLoader[int, int] {
	cfg := DBMappingLoaderConfig{
		TableName:   r.Table,
		KeyColumn:   r.FieldB,
		ValueColumn: r.FieldA,
		Where:       r.getExtConditionSlice(),
	}
	return NewDBBatchSliceLoader[int, int](db, cfg, IntFromRow, IntFromRow)
}

func (r *refPlucker) getExtConditionSlice() (where []string) {
	if s := r.getExtCondition(); s != "" {
		return []string{s}
	}
	return nil
}

// getExtCondition 获取额外的条件
func (r *refPlucker) getExtCondition() string {
	if len(r.ExtKV) == 0 {
		return ""
	}

	if r.extCondition != "" {
		return r.extCondition
	}

	var c []string
	for k, v := range r.ExtKV {
		c = append(c, k+" = '"+v+"'")
	}

	r.extCondition = strings.Join(c, " AND ")
	return r.extCondition
}

// extKVFields 获取额外的字段
func (r *refPlucker) extKVFields() (fields []string, extValue []any) {
	for k, v := range r.ExtKV {
		fields = append(fields, k)
		extValue = append(extValue, v)
	}
	return
}

// table 获取表名
func (r *refPlucker) table(db *gorm.DB) *gorm.DB {
	db = db.Table(r.Table)
	for k, v := range r.ExtKV {
		db = db.Where(k+" = ?", v)
	}
	return db
}

func (r *refPlucker) updateRef(ctx context.Context, db *gorm.DB, a int, b []int, oneField, manyField string) (err error) {
	if err = r.table(db).Where(oneField+" = ?", a).Delete(nil).Error; err != nil {
		return
	}

	if len(b) == 0 {
		return
	}

	keyFields := []string{manyField, oneField}
	extFields, extValue := r.extKVFields()
	if err = gormx.NewWithGormDB(db).BatchInsert(ctx, r.Table, b, func(i int, insert func(...interface{})) {
		insert(append([]interface{}{b[i], a}, extValue...)...)
	}, append(keyFields, extFields...), isql.IsIgnore()); err != nil {
		return
	}
	return
}
