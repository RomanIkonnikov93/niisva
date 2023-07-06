package repository

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/RomanIkonnikov93/niisva/internal/config"
)

type FileRepository struct {
	store string
}

func NewFileRepository(cfg *config.Config) (*FileRepository, error) {

	file, err := resolveFile(cfg.UsersPathStore)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &FileRepository{
		store: cfg.UsersPathStore,
	}, nil
}

func resolveFile(path string) (*os.File, error) {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (r *FileRepository) Add(ctx context.Context, path string) error {

	file, err := resolveFile(r.store)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(path + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) GetAll(ctx context.Context) (users map[string]struct{}, err error) {

	file, err := resolveFile(r.store)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvRead := csv.NewReader(file)
	rec, err := csvRead.ReadAll()
	if err != nil {
		return nil, err
	}

	users = make(map[string]struct{})

	for _, record := range rec {
		if len(record) != 1 {
			continue
		}
		users[record[0]] = struct{}{}
	}

	return users, nil
}
