package book

import (
	"context"

	"go.uber.org/zap"
)

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.bookRepo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete book", zap.String("book_id", id), zap.Error(err))
		return err
	}

	s.log.Info("book deleted successfully", zap.String("book_id", id))
	return nil
}
