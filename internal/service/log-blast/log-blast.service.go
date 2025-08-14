package logblast

import (
	"context"
	"log"
	"waha-job-processing/internal/database/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogBlastService struct {
	DB *pgxpool.Pool
}

func NewLogBlastService(database *pgxpool.Pool) *LogBlastService {
	return &LogBlastService{DB: database}
}

func (s *LogBlastService) CreateLogBlast(logBlastParam repository.CreateLogBlastParams) (*repository.LogBlast, error) {
	ctx := context.Background()
	query := repository.New(s.DB)
	logBlast, err := query.CreateLogBlast(ctx, logBlastParam)
	if err != nil {
		log.Printf("Failed to create log blast: %v", logBlast)
		return nil, err
	}
	log.Printf("Created log blast: %v", logBlast)
	return &logBlast, nil
}

func (s *LogBlastService) UpdateLogBlast(logBlastParam repository.UpdateLogBlastParams) (*repository.LogBlast, error) {
	ctx := context.Background()
	query := repository.New(s.DB)
	logBlast, err := query.UpdateLogBlast(ctx, logBlastParam)
	if err != nil {
		log.Printf("Failed to update log blast: %v", logBlast)
		return nil, err
	}
	log.Printf("Updated log blast: %v", logBlast)
	return &logBlast, nil
}
