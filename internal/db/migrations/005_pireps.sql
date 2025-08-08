
ALTER TABLE va_user_roles 
ADD COLUMN IF NOT EXISTS callsign VARCHAR(40);

ALTER TYPE va_role RENAME VALUE 'airline_manager' TO 'staff';
