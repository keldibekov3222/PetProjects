DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' 
        AND column_name = 'updated_at'
    ) THEN
        ALTER TABLE users ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
    END IF;
END $$; 