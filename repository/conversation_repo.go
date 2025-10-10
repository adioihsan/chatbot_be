package repository

import (
	"context"
	"time"

	"cms-octo-chat-api/model"

	"github.com/google/uuid"
)

func (r *BaseRepository) CreateConversation(ctx context.Context, title string, userID int64) (*model.Conversation, error) {
	c := &model.Conversation{
		Title:  title,
		UserID: userID,
	}
	if c.Title == "" {
		c.Title = "Untitled chat"
	}
	if err := r.DB.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *BaseRepository) GetConversationByPublicID(ctx context.Context, pid uuid.UUID, userID int64) (*model.Conversation, error) {
	var c model.Conversation
	q := r.DB.WithContext(ctx).Where("public_id = ?", pid)
	q = q.Where("user_id = ?", userID)

	if err := q.First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *BaseRepository) ListConversation(ctx context.Context, userID int64, limit int) ([]model.Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var out []model.Conversation
	q := r.DB.WithContext(ctx).Model(&model.Conversation{}).Where("user_id = ?", userID)

	if err := q.Order("updated_at DESC").Limit(limit).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *BaseRepository) RenameConversation(ctx context.Context, pid uuid.UUID, userID int64, title string) error {
	q := r.DB.WithContext(ctx).Model(&model.Conversation{}).Where("public_id = ?", pid)
	q = q.Where("user_id = ?", userID)

	return q.Updates(map[string]any{"title": title}).Error
}

func (r *BaseRepository) DeleteConversation(ctx context.Context, pid uuid.UUID, userID string) error {
	q := r.DB.WithContext(ctx).Where("public_id = ?", pid)
	q = q.Where("user_id = ?", userID)
	return q.Delete(&model.Conversation{}).Error
}

func (r *BaseRepository) TouchConversation(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).
		UpdateColumn("updated_at", time.Now()).Error
}
