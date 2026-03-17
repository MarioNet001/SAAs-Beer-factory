package audit

import (
	"fmt"
	"time"
)

type AuditService interface {
	LogEvent(eventType string, details string) error
}

type auditService struct{}

func NewAuditService() AuditService {
	return &auditService{}
}

func (s *auditService) LogEvent(eventType string, details string) error {
	// Implementation: In real project, save to DB
	fmt.Printf("[%s] Audit Event: %s - %s\n", time.Now().Format(time.RFC3339), eventType, details)
	return nil
}
