#!/bin/bash

# MySQL Setup Script for Betzone
# This script will help you set up MySQL database and user

set -e

echo "================================"
echo "Betzone MySQL Setup Script"
echo "================================"
echo ""

# Check if MySQL is installed
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL is not installed."
    echo ""
    echo "Install MySQL with:"
    echo "  Ubuntu/Debian: sudo apt-get install mysql-server"
    echo "  macOS: brew install mysql"
    echo ""
    exit 1
fi

echo "✓ MySQL is installed"
echo ""

# Try to connect and set up database
echo "Attempting to connect to MySQL..."
echo "You may need to enter your MySQL root password"
echo ""

# Create a temporary SQL file with all commands
TEMP_SQL=$(mktemp)
cat > "$TEMP_SQL" << 'EOF'
-- Create database
CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user if not exists
CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';

-- Grant all privileges
GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';

-- Apply changes
FLUSH PRIVILEGES;

-- Show results
SELECT 'Database and user created successfully!' as message;
SHOW DATABASES LIKE 'betzone';
EOF

echo "Running MySQL setup commands..."
mysql -u root -p < "$TEMP_SQL"

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Database and user created successfully!"
    echo ""
    echo "Your MySQL credentials are:"
    echo "  Host: localhost"
    echo "  Port: 3306"
    echo "  User: betzone"
    echo "  Password: betzone_password_2024"
    echo ""
    echo "Updating .env file..."
    
    # Update .env file
    if [ -f ".env" ]; then
        sed -i "s/DB_USER=.*/DB_USER=betzone/" .env
        sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=betzone_password_2024/" .env
        sed -i "s/DB_NAME=.*/DB_NAME=betzone/" .env
        echo "✓ .env file updated"
    else
        echo "⚠ .env file not found, creating from template..."
        cp .env.example .env
        sed -i "s/your_mysql_password_here/betzone_password_2024/" .env
        echo "✓ .env file created"
    fi
    
    echo ""
    echo "✓ Setup complete!"
    echo ""
    echo "Next steps:"
    echo "  1. Run the application: make dev"
    echo "  2. Test with: curl http://localhost:8080/health"
else
    echo ""
    echo "❌ Failed to set up database"
    echo ""
    echo "Possible solutions:"
    echo "  1. Check your MySQL root password"
    echo "  2. Ensure MySQL is running: sudo systemctl start mysql"
    echo "  3. Try manual setup:"
    echo "     mysql -u root -p"
    echo "     Then paste the commands from MySQL_SETUP_MANUAL.sql"
fi

# Cleanup
rm -f "$TEMP_SQL"
