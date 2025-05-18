-- Alter Users table
ALTER TABLE virtual_airlines ADD COLUMN discord_server_id VARCHAR(32) UNIQUE;
CREATE UNIQUE INDEX idx_va_discord_server_id ON virtual_airlines(discord_server_id);

ALTER TABLE users ADD COLUMN otp VARCHAR(6);