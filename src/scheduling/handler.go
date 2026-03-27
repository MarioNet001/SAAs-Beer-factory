package schedulingapi

import (
	"context"
	"sistema-gestion-beer/src/domain/scheduling"
)

type SchedulingHandler struct {
	service scheduling.SchedulingService
}

func NewHandler(service scheduling.SchedulingService) *SchedulingHandler {
	return &SchedulingHandler{service: service}
}

func (h *SchedulingHandler) HandleCreateSchedule(ctx context.Context, s *scheduling.Schedule) error {
	return h.service.CreateSchedule(ctx, s)
}

func (h *SchedulingHandler) HandleRegisterTank(ctx context.Context, t *scheduling.Tank) error {
	return h.service.RegisterTank(ctx, t)
}

func (h *SchedulingHandler) HandleListSchedulesByTank(ctx context.Context, tankID string) ([]*scheduling.Schedule, error) {
	return h.service.ListSchedulesByTank(ctx, tankID)
}
