
ALTER TABLE va_user_roles 
ADD COLUMN IF NOT EXISTS callsign VARCHAR(40)

ALTER TYPE va_role RENAME VALUE 'airline_manager' TO 'staff';

-- Drop the old unique constraint
ALTER TABLE va_user_roles
DROP CONSTRAINT va_user_roles_user_id_va_id_role_key;

-- Add the new unique constraint
ALTER TABLE va_user_roles
ADD CONSTRAINT va_user_roles_user_id_va_id_key UNIQUE (user_id, va_id);
