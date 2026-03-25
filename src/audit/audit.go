package audit

import (
	"database/sql"
	"fmt"
	"time"
)

type AuditService interface {
	LogEvent(eventType string, details string) error
}

type auditService struct {
	db *sql.DB
}

func NewAuditService(db *sql.DB) AuditService {
	return &auditService{db: db}
}

func (s *auditService) LogEvent(eventType string, details string) error {
	_, err := s.db.Exec("INSERT INTO audit_logs (event_type, details, created_at) VALUES ($1, $2, $3)", eventType, details, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log audit event: %w", err)
	}
	return nil
}
