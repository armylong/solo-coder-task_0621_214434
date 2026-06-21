package notification

import (
	"time"

	"github.com/armylong/go-library/service/sqlite"
)

const (
	AnnouncementStatusDraft     = 0
	AnnouncementStatusPublished = 1
	AnnouncementStatusArchived  = 2

	AnnouncementTypeSystem      = "system"
	AnnouncementTypeMaintenance = "maintenance"
	AnnouncementTypeActivity    = "activity"
	AnnouncementTypeUpdate      = "update"

	AnnouncementPriorityLow    = 0
	AnnouncementPriorityMedium = 1
	AnnouncementPriorityHigh   = 2
	AnnouncementPriorityUrgent = 3

	AnnouncementTargetAll       = "all"
	AnnouncementTargetAdmin     = "admin"
	AnnouncementTargetSpecific  = "specific"
)

type TbAnnouncement struct {
	ID             int64      `json:"id" db:"pk"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	Priority       int        `json:"priority"`
	TargetAudience string     `json:"target_audience"`
	TargetUids     string     `json:"target_uids"`
	Status         int        `json:"status"`
	CreatedBy      int64      `json:"created_by"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type tbAnnouncementModel struct{}

var TbAnnouncementModel = &tbAnnouncementModel{}

func init() {
	_ = TbAnnouncementModel.CreateTable()
}

func (m *tbAnnouncementModel) TableName() string {
	return "tb_announcement"
}

func (m *tbAnnouncementModel) CreateTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS tb_announcement (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		type TEXT NOT NULL DEFAULT 'system',
		priority INTEGER NOT NULL DEFAULT 0,
		target_audience TEXT NOT NULL DEFAULT 'all',
		target_uids TEXT NOT NULL DEFAULT '',
		status INTEGER NOT NULL DEFAULT 0,
		created_by INTEGER NOT NULL DEFAULT 0,
		published_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := sqlite.DB.DB().Exec(sql)
	if err != nil {
		return err
	}
	return sqlite.DB.AutoMigrate(m.TableName(), &TbAnnouncement{})
}

func (m *tbAnnouncementModel) Create(announcement *TbAnnouncement) (int64, error) {
	return sqlite.DB.Insert(m.TableName(), announcement)
}

func (m *tbAnnouncementModel) GetByID(id int64) (*TbAnnouncement, error) {
	var a TbAnnouncement
	a.ID = id
	err := sqlite.DB.GetByPkId(m.TableName(), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (m *tbAnnouncementModel) Update(announcement *TbAnnouncement) error {
	return sqlite.DB.UpdateByPkId(m.TableName(), announcement)
}

func (m *tbAnnouncementModel) Delete(id int64) error {
	a := &TbAnnouncement{ID: id}
	return sqlite.DB.DeleteByPkId(m.TableName(), a)
}

func (m *tbAnnouncementModel) List(status int, page, pageSize int) ([]*TbAnnouncement, int64, error) {
	var list []*TbAnnouncement
	var total int64
	var err error

	if status >= 0 {
		total, err = sqlite.DB.Count(m.TableName(), "status = ?", status)
		if err != nil {
			return nil, 0, err
		}
		offset := (page - 1) * pageSize
		err = sqlite.DB.Find(m.TableName(), &list, "status = ? ORDER BY priority DESC, created_at DESC LIMIT ? OFFSET ?", status, pageSize, offset)
	} else {
		total, err = sqlite.DB.CountAll(m.TableName())
		if err != nil {
			return nil, 0, err
		}
		offset := (page - 1) * pageSize
		err = sqlite.DB.Find(m.TableName(), &list, "1=1 ORDER BY priority DESC, created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *tbAnnouncementModel) ListPublished(targetAudience string, page, pageSize int) ([]*TbAnnouncement, int64, error) {
	var list []*TbAnnouncement
	var total int64
	var err error

	offset := (page - 1) * pageSize

	if targetAudience == "" || targetAudience == AnnouncementTargetAll {
		total, err = sqlite.DB.Count(m.TableName(), "status = ?", AnnouncementStatusPublished)
		if err != nil {
			return nil, 0, err
		}
		err = sqlite.DB.Find(m.TableName(), &list, "status = ? ORDER BY priority DESC, published_at DESC LIMIT ? OFFSET ?", AnnouncementStatusPublished, pageSize, offset)
	} else {
		total, err = sqlite.DB.Count(m.TableName(), "status = ? AND (target_audience = ? OR target_audience = ?)", AnnouncementStatusPublished, AnnouncementTargetAll, targetAudience)
		if err != nil {
			return nil, 0, err
		}
		err = sqlite.DB.Find(m.TableName(), &list, "status = ? AND (target_audience = ? OR target_audience = ?) ORDER BY priority DESC, published_at DESC LIMIT ? OFFSET ?", AnnouncementStatusPublished, AnnouncementTargetAll, targetAudience, pageSize, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *tbAnnouncementModel) GetPublishedByIDs(ids []int64) ([]*TbAnnouncement, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*TbAnnouncement
	placeholders := ""
	args := make([]interface{}, 0, len(ids)+1)
	args = append(args, AnnouncementStatusPublished)
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	err := sqlite.DB.Find(m.TableName(), &list, "status = ? AND id IN ("+placeholders+")", args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}
