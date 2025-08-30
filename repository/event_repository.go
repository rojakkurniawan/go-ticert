package repository

import (
	"ticert/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(event *entity.Event) error
	GetEventByID(id uuid.UUID) (*entity.Event, error)
	GetEventByTitle(title string) (*entity.Event, error)
	GetEvents(page, limit int, search string, orderBy string) ([]*entity.Event, int64, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(id uuid.UUID) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(event *entity.Event) error {
	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetEventByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetEventByTitle(title string) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("title = ?", title).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetEvents(page, limit int, search string, orderBy string) ([]*entity.Event, int64, error) {
	var events []*entity.Event
	var total int64

	query := r.db.Model(&entity.Event{})

	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR organizer LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if orderBy != "" {
		validOrderings := map[string]string{
			"asc":       "created_at ASC",
			"desc":      "created_at DESC",
			"date_asc":  "start_date ASC",
			"date_desc": "start_date DESC",
			"time_asc":  "start_time ASC",
			"time_desc": "start_time DESC",
		}

		validOrderBy := ""
		if v, exists := validOrderings[orderBy]; exists {
			validOrderBy = v
		}

		if validOrderBy != "" {
			query = query.Order(validOrderBy)
		}

	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Categories").Offset((page - 1) * limit).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepository) UpdateEvent(event *entity.Event) error {
	if err := r.db.Model(&entity.Event{}).Where("id = ?", event.ID).Updates(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) DeleteEvent(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Where("event_id = ?", id).Delete(&entity.Category{}).Error; err != nil {
			return err
		}

		if err := r.db.Delete(&entity.Event{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
