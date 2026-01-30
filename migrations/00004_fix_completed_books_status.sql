-- +goose Up
-- Add 'on_hold' status to books status check constraint
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_status_check;
ALTER TABLE books ADD CONSTRAINT books_status_check 
    CHECK (status IN ('available', 'reading', 'reserved', 'requested', 'on_hold'));

-- Fix books that have completed reading history but are still in "reading" status
-- These should be "on_hold" status
UPDATE books 
SET status = 'on_hold', updated_at = NOW()
WHERE id IN (
    SELECT DISTINCT b.id 
    FROM books b
    INNER JOIN reading_history rh ON b.id = rh.book_id
    WHERE b.status = 'reading'
    AND rh.is_completed = true
    AND rh.end_date IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM reading_history rh2 
        WHERE rh2.book_id = b.id 
        AND rh2.end_date IS NULL
    )
);

-- +goose Down
-- Remove 'on_hold' status from books status check constraint
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_status_check;
ALTER TABLE books ADD CONSTRAINT books_status_check 
    CHECK (status IN ('available', 'reading', 'reserved', 'requested'));
