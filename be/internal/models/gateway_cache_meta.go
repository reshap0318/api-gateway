package models

// GatewayCacheMeta is a singleton row (id=1) counting Service/Route CUD — lets
// RouteManager's periodic refresh skip rebuilding when nothing changed. A monotonic
// counter, not a timestamp, so it stays correct regardless of clock sync across instances.
type GatewayCacheMeta struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Version uint64 `gorm:"column:version;not null" json:"version"`
}

// TableName specifies the table name for GatewayCacheMeta model.
func (GatewayCacheMeta) TableName() string {
	return "gateway_cache_meta"
}
