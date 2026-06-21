package notification

import (
	"time"

	"github.com/armylong/go-library/service/sqlite"
)

type TbNotificationPreference struct {
	ID               int64     `json:"id" db:"pk"`
	Uid              int64     `json:"uid"`
	NotificationType string    `json:"notification_type"`
	Enabled          int       `json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type tbNotificationPreferenceModel struct{}

var TbNotificationPreferenceModel = &tbNotificationPreferenceModel{}

func init() {
	_ = TbNotificationPreferenceModel.CreateTable()
}

func (m *tbNotificationPreferenceModel) TableName() string {
	return "tb_notification_preference"
}

func (m *tbNotificationPreferenceModel) CreateTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS tb_notification_preference (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid INTEGER NOT NULL,
		notification_type TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(uid, notification_type)
	)
	`
	_, err := sqlite.DB.DB().Exec(sql)
	if err != nil {
		return err
	}
	return sqlite.DB.AutoMigrate(m.TableName(), &TbNotificationPreference{})
}

func (m *tbNotificationPreferenceModel) GetByUid(uid int64) ([]*TbNotificationPreference, error) {
	var list []*TbNotificationPreference
	err := sqlite.DB.Find(m.TableName(), &list, "uid = ?", uid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *tbNotificationPreferenceModel) GetByUidAndType(uid int64, notificationType string) (*TbNotificationPreference, error) {
	var pref TbNotificationPreference
	err := sqlite.DB.FindOne(m.TableName(), &pref, "uid = ? AND notification_type = ?", uid, notificationType)
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func (m *tbNotificationPreferenceModel) Upsert(pref *TbNotificationPreference) error {
	existing, err := m.GetByUidAndType(pref.Uid, pref.NotificationType)
	if err == nil && existing != nil {
		existing.Enabled = pref.Enabled
		return sqlite.DB.UpdateByPkId(m.TableName(), existing)
	}
	_, err = sqlite.DB.Insert(m.TableName(), pref)
	return err
}

func (m *tbNotificationPreferenceModel) IsEnabled(uid int64, notificationType string) bool {
	pref, err := m.GetByUidAndType(uid, notificationType)
	if err != nil || pref == nil {
		return true
	}
	return pref.Enabled == 1
}

func (m *tbNotificationPreferenceModel) BatchUpsert(preferences []*TbNotificationPreference) error {
	for _, pref := range preferences {
		if err := m.Upsert(pref); err != nil {
			return err
		}
	}
	return nil
}
