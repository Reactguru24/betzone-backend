#!/bin/bash

# Interactive MySQL Setup for Betzone
# This script will help you create the database if you have the root password

set -e

echo "╔════════════════════════════════════════════════╗"
echo "║   Betzone MySQL Interactive Setup              ║"
echo "╚════════════════════════════════════════════════╝"
echo ""

# Test if MySQL is accessible without password
echo "Testing MySQL connection..."

if sudo mysql -u root -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✓ MySQL root user has no password"
    echo ""
    echo "Setting up database..."
    
    # Create database and user
    sudo mysql -u root << EOF
CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';
GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';
FLUSH PRIVILEGES;
EOF
    
    if [ $? -eq 0 ]; then
        echo "✓ Database setup successful!"
    fi
    
elif mysql -u root -p -e "SELECT 1;" > /dev/null 2>&1; then
    # Interactive password prompt
    echo ""
    echo "MySQL root user requires a password."
    echo "You will be prompted for your MySQL root password."
    echo ""
    echo "Setting up database..."
    
    mysql -u root -p << 'EOF'
CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';
GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';
FLUSH PRIVILEGES;
SHOW DATABASES;
EOF
    
    echo ""
    echo "✓ Database setup complete!"
    
else
    echo "❌ Cannot connect to MySQL"
    echo ""
    echo "Please try manual setup:"
    echo ""
    echo "1. Open terminal and run:"
    echo "   mysql -u root -p"
    echo ""
    echo "2. When prompted, paste these commands:"
    echo "   CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4;"
    echo "   CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';"
    echo "   GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';"
    echo "   FLUSH PRIVILEGES;"
    echo "   EXIT;"
    exit 1
fi

echo ""
echo "=============================================="
echo "Verifying setup..."
echo "=============================================="
echo ""

# Verify the database was created
if mysql -u betzone -p'betzone_password_2024' -h localhost -e "USE betzone; SHOW TABLES;" > /dev/null 2>&1; then
    echo "✓ Database 'betzone' created successfully"
    echo "✓ User 'betzone' created successfully"
    echo "✓ User has proper permissions"
    echo ""
    echo "✓ All checks passed!"
    echo ""
    echo "Your .env file is configured with:"
    echo "  DB_USER=betzone"
    echo "  DB_PASSWORD=betzone_password_2024"
    echo "  DB_NAME=betzone"
    echo ""
    echo "You can now run: make dev"
else
    echo "⚠ Could not verify database"
fi
