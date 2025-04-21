-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discord_id VARCHAR(32) UNIQUE NOT NULL,
    if_community_id VARCHAR(10),
    if_api_id UUID,
    is_active BOOLEAN DEFAULT false,
    username TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create virtual airlines table
CREATE TABLE virtual_airlines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create enum for VA roles
DO $$ BEGIN
    CREATE TYPE va_role AS ENUM ('pilot', 'airline_manager', 'admin');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create VA user roles table
CREATE TABLE va_user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    va_id UUID REFERENCES virtual_airlines(id),
    role va_role NOT NULL,
    is_active BOOLEAN DEFAULT true,
    joined_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, va_id, role)
);


-- =====================
-- Indexes for performance
-- =====================

-- users
CREATE UNIQUE INDEX idx_users_discord_id ON users(discord_id);
CREATE INDEX idx_users_username ON users(username);

-- va_user_roles
CREATE INDEX idx_va_user_roles_user_id ON va_user_roles(user_id);
CREATE INDEX idx_va_user_roles_va_id ON va_user_roles(va_id);
CREATE INDEX idx_va_user_roles_user_va ON va_user_roles(user_id, va_id);
