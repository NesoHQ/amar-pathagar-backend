-- +goose Up
-- Remove the old unique constraint and add a partial unique index
-- This allows multiple approved/rejected requests but only one pending request per book per user

-- Drop the old constraint
ALTER TABLE book_requests DROP CONSTRAINT IF EXISTS book_requests_book_id_user_id_status_key;

-- Create partial unique index: only one pending request per book per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_book_requests_pending_unique 
ON book_requests(book_id, user_id) 
WHERE status = 'pending';

-- +goose Down
-- Revert to the old constraint
DROP INDEX IF EXISTS idx_book_requests_pending_unique;
ALTER TABLE book_requests ADD CONSTRAINT book_requests_book_id_user_id_status_key UNIQUE(book_id, user_id, status);
