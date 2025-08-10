package logblast

import (
	"context"
	"log"
	"waha-job-processing/internal/database/db"
	"waha-job-processing/internal/database/repository"
)

func CreateLogBlast(logBlastParam repository.CreateLogBlastParams) (*repository.LogBlast, error) {
	ctx := context.Background()
	conn := db.New(ctx)
	defer conn.Close(ctx)

	query := repository.New(conn)

	logBlast, err := query.CreateLogBlast(ctx, logBlastParam)
	if err != nil {
		log.Printf("Failed to create log blast: %v", logBlast)
		return nil, err
	}

	log.Printf("Created log blast: %v", logBlast)
	return &logBlast, nil
}

func UpdateLogBlast(logBlastParam repository.UpdateLogBlastParams) (*repository.LogBlast, error) {
	ctx := context.Background()
	conn := db.New(ctx)
	defer conn.Close(ctx)

	query := repository.New(conn)

	logBlast, err := query.UpdateLogBlast(ctx, logBlastParam)
	if err != nil {
		log.Printf("Failed to update log blast: %v", logBlast)
		return nil, err
	}

	log.Printf("Updated log blast: %v", logBlast)
	return &logBlast, nil
}
