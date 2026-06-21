package notification

import (
	"time"

	"github.com/armylong/go-library/service/sqlite"
)

type TbNotification struct {
	ID             int64      `json:"id" db:"pk"`
	Uid            int64      `json:"uid"`
	AnnouncementID int64      `json:"announcement_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	IsRead         int        `json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type tbNotificationModel struct{}

var TbNotificationModel = &tbNotificationModel{}

func init() {
	_ = TbNotificationModel.CreateTable()
}

func (m *tbNotificationModel) TableName() string {
	return "tb_notification"
}

func (m *tbNotificationModel) CreateTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS tb_notification (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid INTEGER NOT NULL,
		announcement_id INTEGER NOT NULL DEFAULT 0,
		title TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		type TEXT NOT NULL DEFAULT 'system',
		is_read INTEGER NOT NULL DEFAULT 0,
		read_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := sqlite.DB.DB().Exec(sql)
	if err != nil {
		return err
	}
	return sqlite.DB.AutoMigrate(m.TableName(), &TbNotification{})
}

func (m *tbNotificationModel) Create(notification *TbNotification) (int64, error) {
	return sqlite.DB.Insert(m.TableName(), notification)
}

func (m *tbNotificationModel) CreateBatch(notifications []*TbNotification) error {
	for _, n := range notifications {
		_, err := sqlite.DB.Insert(m.TableName(), n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *tbNotificationModel) GetByID(id int64) (*TbNotification, error) {
	var n TbNotification
	n.ID = id
	err := sqlite.DB.GetByPkId(m.TableName(), &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (m *tbNotificationModel) ListByUid(uid int64, notificationType string, page, pageSize int) ([]*TbNotification, int64, error) {
	var list []*TbNotification
	var total int64
	var err error
	offset := (page - 1) * pageSize

	if notificationType != "" {
		total, err = sqlite.DB.Count(m.TableName(), "uid = ? AND type = ?", uid, notificationType)
		if err != nil {
			return nil, 0, err
		}
		err = sqlite.DB.Find(m.TableName(), &list, "uid = ? AND type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?", uid, notificationType, pageSize, offset)
	} else {
		total, err = sqlite.DB.Count(m.TableName(), "uid = ?", uid)
		if err != nil {
			return nil, 0, err
		}
		err = sqlite.DB.Find(m.TableName(), &list, "uid = ? ORDER BY created_at DESC LIMIT ? OFFSET ?", uid, pageSize, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *tbNotificationModel) MarkAsRead(id int64) error {
	now := time.Now()
	_, err := sqlite.DB.DB().Exec("UPDATE tb_notification SET is_read = 1, read_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", now, id)
	return err
}

func (m *tbNotificationModel) MarkAsReadByUid(uid int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	placeholders := ""
	args := []interface{}{now, uid}
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	_, err := sqlite.DB.DB().Exec("UPDATE tb_notification SET is_read = 1, read_at = ?, updated_at = CURRENT_TIMESTAMP WHERE uid = ? AND id IN ("+placeholders+")", args...)
	return err
}

func (m *tbNotificationModel) MarkAllAsRead(uid int64) error {
	now := time.Now()
	_, err := sqlite.DB.DB().Exec("UPDATE tb_notification SET is_read = 1, read_at = ?, updated_at = CURRENT_TIMESTAMP WHERE uid = ? AND is_read = 0", now, uid)
	return err
}

func (m *tbNotificationModel) DeleteReadByUid(uid int64) error {
	return sqlite.DB.DeleteByWhere(m.TableName(), "uid = ? AND is_read = 1", uid)
}

func (m *tbNotificationModel) DeleteByUid(uid int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := ""
	args := []interface{}{uid}
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	_, err := sqlite.DB.DB().Exec("DELETE FROM tb_notification WHERE uid = ? AND id IN ("+placeholders+")", args...)
	return err
}

func (m *tbNotificationModel) CountUnreadByUid(uid int64) (int64, error) {
	return sqlite.DB.Count(m.TableName(), "uid = ? AND is_read = 0", uid)
}

func (m *tbNotificationModel) CountUnreadByUidGroupByType(uid int64) ([]map[string]interface{}, error) {
	rows, err := sqlite.DB.DB().Query("SELECT type, COUNT(*) as count FROM tb_notification WHERE uid = ? AND is_read = 0 GROUP BY type", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var nType string
		var count int64
		if err := rows.Scan(&nType, &count); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"type":  nType,
			"count": count,
		})
	}
	return result, nil
}

func (m *tbNotificationModel) ExistsByUidAndAnnouncementID(uid int64, announcementID int64) (bool, error) {
	count, err := sqlite.DB.Count(m.TableName(), "uid = ? AND announcement_id = ?", uid, announcementID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
