-- Betzone MySQL Setup Commands
-- Paste these commands inside the MySQL prompt
-- 
-- First connect with: mysql -u root -p
-- Then paste the commands below

-- Step 1: Create the database with UTF8MB4 support
CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Step 2: Create a dedicated user for the application
-- Password: betzone_password_2024 (change this to something stronger in production!)
CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';

-- Step 3: Grant all privileges on the betzone database
GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';

-- Step 4: Apply the privileges
FLUSH PRIVILEGES;

-- Step 5: Verify the setup
-- You should see 'betzone' database in the list below
SHOW DATABASES;

-- Step 6: Verify the user was created
SELECT user, host FROM mysql.user WHERE user='betzone';

-- Step 7: Switch to the betzone database
USE betzone;

-- Step 8: Show tables (should be empty initially, tables will be created by the app)
SHOW TABLES;

-- All done! You can now exit MySQL with:
-- EXIT;
