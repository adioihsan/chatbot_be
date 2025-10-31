package repository

import (
	"context"
	"fmt"

	"cms-octo-chat-api/model"

	"github.com/google/uuid"
)

func (r *BaseRepository) ListMessageByConversationPID(ctx context.Context, userID int64, convoPID uuid.UUID, limit int, beforePublicID *uuid.UUID) ([]model.Message, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	var conv model.Conversation
	if err := r.DB.WithContext(ctx).Select("id").Where("user_id = ?", userID).First(&conv, "public_id = ?", convoPID).Error; err != nil {
		return nil, err
	}

	q := r.DB.WithContext(ctx).Model(&model.Message{}).Where("conversation_id = ?", conv.ID)
	if beforePublicID != nil {
		var before model.Message
		if err := r.DB.WithContext(ctx).Select("created_at").First(&before, "public_id = ?", *beforePublicID).Error; err == nil {
			q = q.Where("created_at < ?", before.CreatedAt)
		}
	}

	subQuery := q.Order("created_at DESC").Limit(limit)

	var out []model.Message
	if err := r.DB.Table("(?) as main_data", subQuery).Order("main_data.created_at ASC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *BaseRepository) AppendMessage(ctx context.Context, m *model.Message) (*model.Message, error) {
	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *BaseRepository) RebuildMessageContentFTS(ctx context.Context, messageId int64) error {
	baseExpr := `setweight(to_tsvector('pg_catalog.english', unaccent(coalesce(content,''))), 'B')`
	sql := fmt.Sprintf("UPDATE messages SET content_vector = %s", baseExpr)
	sql += fmt.Sprintf(" WHERE id = %d", messageId)

	return r.DB.WithContext(ctx).Exec(sql).Error
}
