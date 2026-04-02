package gormtil

import (
	"context"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

// KeyFromRow 从查询行中解析出 key 的函数，用于泛型 Loader
type KeyFromRow[K comparable] func(rowKey any) (K, error)

// ValueFromRow 从查询行中解析出 value 的函数，用于泛型 Loader
type ValueFromRow[V any] func(rowValue any) (V, error)

// DBMappingLoaderConfig 通用 DB 映射 Loader 的配置：表名、键列、值列、可选 where
type DBMappingLoaderConfig struct {
	TableName   string   // 表名
	KeyColumn   string   // 作为 key 的列名，如 "id"、"process_id"
	ValueColumn string   // 作为 value 的列/表达式，如 "name"、"count(*)"
	Where       []string // 额外 WHERE 条件，如 "deleted_at IS NULL"
}

// DBMappingLoader 泛型：按 K 批量查表并得到 map[K]any，填 slice 时用 castPtr 写回
type DBMappingLoader[K comparable] struct {
	loader     *dataloader.Loader[K, any]
	keyFromRow KeyFromRow[K]
}

// NewDBMappingLoader 创建通用 DB 映射 Loader
// cfg: 表名、键列、值列、where；keyFromRow: 把查询结果的 _key 列转成 K（如 IntKeyFromRow、StringKeyFromRow）
func NewDBMappingLoader[K comparable](db *gorm.DB, cfg DBMappingLoaderConfig, keyFromRow KeyFromRow[K]) *DBMappingLoader[K] {
	return &DBMappingLoader[K]{
		loader: dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[any] {
			return batchLoadDBMapping(ctx, db, cfg, keyFromRow, keys)
		}),
		keyFromRow: keyFromRow,
	}
}

// batchLoadDBMapping 批量查询：SELECT keyCol AS _key, valueCol AS _value FROM table WHERE keyCol IN (?)
//
//nolint:lll
func batchLoadDBMapping[K comparable](ctx context.Context, db *gorm.DB, cfg DBMappingLoaderConfig, keyFromRow KeyFromRow[K], keys []K) []*dataloader.Result[any] {
	if len(keys) == 0 {
		return []*dataloader.Result[any]{}
	}

	var rows []map[string]any
	query := db.WithContext(ctx).Table(cfg.TableName)

	for _, w := range cfg.Where {
		if w != "" {
			query = query.Where(w)
		}
	}

	// WHERE key IN (?)
	switch len(keys) {
	case 0:
		return []*dataloader.Result[any]{}
	case 1:
		query = query.Where(cfg.KeyColumn+" = ?", keys[0])
	default:
		query = query.Where(cfg.KeyColumn+" IN ?", keys)
	}

	if isAggregateExpr(cfg.ValueColumn) {
		query = query.Group(cfg.KeyColumn)
	}

	sel := cfg.KeyColumn + " AS _key," + cfg.ValueColumn + " AS _value"
	if err := query.Select(sel).Find(&rows).Error; err != nil {
		results := make([]*dataloader.Result[interface{}], len(keys))
		for i := range keys {
			results[i] = &dataloader.Result[interface{}]{Error: err}
		}
		return results
	}

	// 构建 map[K]interface{}
	valueMap := make(map[K]interface{}, len(rows))
	for _, row := range rows {
		k, err := keyFromRow(row["_key"])
		if err != nil {
			continue
		}
		valueMap[k] = row["_value"]
	}

	// 按 keys 顺序构造 Result
	results := make([]*dataloader.Result[interface{}], len(keys))
	for i, key := range keys {
		if v, ok := valueMap[key]; ok {
			results[i] = &dataloader.Result[interface{}]{Data: v}
		} else {
			results[i] = &dataloader.Result[interface{}]{
				Error: fmt.Errorf("%s key %v not found", cfg.TableName, key),
			}
		}
	}
	return results
}

// isAggregateExpr 判断 valueCol 是否为聚合表达式（如 count(*)、count(distinct x)）
// 简单匹配, 后续根据业务调整
func isAggregateExpr(valueCol string) bool {
	return strings.Contains(valueCol, "(") && strings.Contains(valueCol, ")")
}

// LoadMany 批量查询，返回 map[K]any
func (l *DBMappingLoader[K]) LoadMany(ctx context.Context, keys []K) (map[K]any, error) {
	if len(keys) == 0 {
		return make(map[K]any), nil
	}
	thunks := make([]func() (any, error), len(keys))
	for i, id := range keys {
		// 创建 thunk
		thunks[i] = l.loader.Load(ctx, id)
	}
	resultMap := make(map[K]any, len(keys))
	// 批量执行 thunk
	for i, thunk := range thunks {
		v, err := thunk()
		if err != nil {
			return nil, fmt.Errorf("load key %v: %w", keys[i], err)
		}
		// 写回结果
		resultMap[keys[i]] = v
	}
	return resultMap, nil
}

// DBBatchSliceLoader 泛型：按 K 批量查表并得到 map[K][]V，同一 key 多行聚合成 slice，内部用 DataLoader 批处理
type DBBatchSliceLoader[K comparable, V any] struct {
	loader *dataloader.Loader[K, []V]
}

// NewDBBatchSliceLoader 创建通用「key → []value」批量 Loader
//
//nolint:lll
func NewDBBatchSliceLoader[K comparable, V any](db *gorm.DB, cfg DBMappingLoaderConfig, keyFromRow KeyFromRow[K], valueFromRow ValueFromRow[V]) *DBBatchSliceLoader[K, V] {
	return &DBBatchSliceLoader[K, V]{
		loader: dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[[]V] {
			return batchLoadDBBatchSlice(ctx, db, cfg, keyFromRow, valueFromRow, keys)
		}),
	}
}

// LoadMany 批量查，返回 map[K][]V
func (l *DBBatchSliceLoader[K, V]) LoadMany(ctx context.Context, keys []K) (map[K][]V, error) {
	if len(keys) == 0 {
		return make(map[K][]V), nil
	}
	thunks := make([]func() ([]V, error), len(keys))
	for i, k := range keys {
		thunks[i] = l.loader.Load(ctx, k)
	}
	out := make(map[K][]V, len(keys))
	for i, thunk := range thunks {
		s, err := thunk()
		if err != nil {
			return nil, fmt.Errorf("load key %v: %w", keys[i], err)
		}
		out[keys[i]] = s
	}
	return out, nil
}

// batchLoadDBBatchSlice 批量查询：SELECT keyCol AS _key, valueCol AS _value FROM table WHERE keyCol IN (?)
//
//nolint:lll
func batchLoadDBBatchSlice[K comparable, V any](ctx context.Context, db *gorm.DB, cfg DBMappingLoaderConfig, keyFromRow KeyFromRow[K], valueFromRow ValueFromRow[V], keys []K) []*dataloader.Result[[]V] {
	if len(keys) == 0 {
		return []*dataloader.Result[[]V]{}
	}

	var rows []map[string]any

	// WHERE key IN (?)
	query := db.WithContext(ctx).Table(cfg.TableName)
	for _, w := range cfg.Where {
		if w != "" {
			query = query.Where(w)
		}
	}

	if len(keys) == 1 {
		query = query.Where(cfg.KeyColumn+" = ?", keys[0])
	} else {
		query = query.Where(cfg.KeyColumn+" IN ?", keys)
	}

	sel := cfg.KeyColumn + " AS _key," + cfg.ValueColumn + " AS _value"
	if err := query.Select(sel).Find(&rows).Error; err != nil {
		results := make([]*dataloader.Result[[]V], len(keys))
		for i := range keys {
			results[i] = &dataloader.Result[[]V]{Error: err}
		}
		return results
	}

	m := make(map[K][]V)
	for _, row := range rows {
		k, err := keyFromRow(row["_key"])
		if err != nil {
			continue
		}
		v, err := valueFromRow(row["_value"])
		if err != nil {
			continue
		}
		m[k] = append(m[k], v)
	}

	results := make([]*dataloader.Result[[]V], len(keys))
	for i, key := range keys {
		if s, ok := m[key]; ok {
			results[i] = &dataloader.Result[[]V]{Data: s}
		} else {
			results[i] = &dataloader.Result[[]V]{Data: nil}
		}
	}
	return results
}
