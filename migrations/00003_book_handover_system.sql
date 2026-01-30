-- +goose Up
-- Book handover system with reading periods and delivery tracking

-- Add reading period fields to reading_history
ALTER TABLE reading_history 
ADD COLUMN IF NOT EXISTS due_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS is_completed BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS next_reader_id UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS delivery_status VARCHAR(20) DEFAULT 'not_started' CHECK (delivery_status IN ('not_started', 'in_transit', 'delivered')),
ADD COLUMN IF NOT EXISTS marked_delivered_at TIMESTAMP;

-- Book handover threads table
CREATE TABLE IF NOT EXISTS handover_threads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    current_holder_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    next_holder_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reading_history_id UUID REFERENCES reading_history(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
    handover_due_date TIMESTAMP NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Handover messages table
CREATE TABLE IF NOT EXISTS handover_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    thread_id UUID NOT NULL REFERENCES handover_threads(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_system_message BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add max reading period to books (in days, default 14 days)
ALTER TABLE books 
ADD COLUMN IF NOT EXISTS max_reading_days INTEGER DEFAULT 14;

-- Indexes for performance
CREATE INDEX idx_reading_history_due_date ON reading_history(due_date);
CREATE INDEX idx_reading_history_next_reader ON reading_history(next_reader_id);
CREATE INDEX idx_reading_history_delivery_status ON reading_history(delivery_status);
CREATE INDEX idx_handover_threads_book ON handover_threads(book_id);
CREATE INDEX idx_handover_threads_current_holder ON handover_threads(current_holder_id);
CREATE INDEX idx_handover_threads_next_holder ON handover_threads(next_holder_id);
CREATE INDEX idx_handover_threads_status ON handover_threads(status);
CREATE INDEX idx_handover_messages_thread ON handover_messages(thread_id);
CREATE INDEX idx_handover_messages_user ON handover_messages(user_id);

-- Trigger for handover_threads updated_at
CREATE TRIGGER update_handover_threads_updated_at BEFORE UPDATE ON handover_threads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS handover_messages CASCADE;
DROP TABLE IF EXISTS handover_threads CASCADE;

ALTER TABLE reading_history 
DROP COLUMN IF EXISTS due_date,
DROP COLUMN IF EXISTS is_completed,
DROP COLUMN IF EXISTS completed_at,
DROP COLUMN IF EXISTS next_reader_id,
DROP COLUMN IF EXISTS delivery_status,
DROP COLUMN IF EXISTS marked_delivered_at;

ALTER TABLE books 
DROP COLUMN IF EXISTS max_reading_days;
