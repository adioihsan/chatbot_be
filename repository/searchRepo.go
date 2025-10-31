package repository

import (
	"cms-octo-chat-api/model"
	"context"
	"fmt"
)

func (r *BaseRepository) FindChats(ctx context.Context, searchQuery string, limit int, userID int64) ([]model.SearchResult, error) {

	tsquerySQL := "websearch_to_tsquery('pg_catalog.english', ?)"

	sql := fmt.Sprintf(`
		WITH q AS (
		  SELECT %s AS query
		)
		SELECT DISTINCT ON (c.public_id)
		  c.public_id as conversation_pid,
		  c.title as title,
		  ts_rank_cd(m.content_vector, q.query) AS rank,
		  ts_headline('pg_catalog.english', m.content, q.query,
		    'MaxFragments=2, MinWords=2, MaxWords=10, ShortWord=2, HighlightAll=FALSE') AS content
		FROM conversations c join messages m ON c.id = m.conversation_id, q
		WHERE m.content_vector @@ q.query
		ORDER BY c.public_id, rank DESC, c.updated_at DESC
		LIMIT ?;
	`, tsquerySQL)

	rows := make([]model.SearchResult, 0, limit)
	if err := r.DB.WithContext(ctx).Raw(sql, searchQuery, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil

}
