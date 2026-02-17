#!/bin/bash

# MySQL Root Password Helper
# This script will help you identify and fix the root password issue

echo "╔═══════════════════════════════════════════════════════╗"
echo "║  MySQL Root Password Troubleshooter                   ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo ""

echo "Let's find your MySQL root password..."
echo ""

# Test 1: Try without password
echo "Test 1: Trying connection without password..."
if mysql -u root -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✓ MySQL root has NO password!"
    echo ""
    echo "Your MySQL root user doesn't require a password."
    echo "This means you can run:"
    echo ""
    echo "  mysql -u root"
    echo ""
    echo "To set up the database, run:"
    echo "  bash final_setup.sh"
    exit 0
fi

echo "✗ MySQL root requires a password"
echo ""

# Test 2: Try common default passwords
echo "Test 2: Trying common default passwords..."
echo ""

PASSWORDS=("root" "password" "12345" "123456" "admin" "mysql" "")

for pass in "${PASSWORDS[@]}"; do
    if [ -z "$pass" ]; then
        if mysql -u root -e "SELECT 1;" > /dev/null 2>&1; then
            echo "✓ Found it! MySQL root password is: [empty/no password]"
            exit 0
        fi
    else
        if mysql -u root -p"$pass" -e "SELECT 1;" > /dev/null 2>&1; then
            echo "✓ Found it! MySQL root password is: $pass"
            echo ""
            echo "Save this for later!"
            exit 0
        fi
    fi
done

echo "✗ Common passwords didn't work"
echo ""

# If we get here, they need to provide the password manually
echo "╔═══════════════════════════════════════════════════════╗"
echo "║  Please Provide Your MySQL Root Password              ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo ""
echo "I couldn't find your MySQL root password automatically."
echo ""
echo "You have 3 options:"
echo ""
echo "OPTION 1: Provide the password here"
echo "   (Type your password and press Enter)"
echo ""
read -sp "MySQL root password: " MYSQL_PASSWORD
echo ""
echo ""

# Test the provided password
if mysql -u root -p"$MYSQL_PASSWORD" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✓ Password confirmed!"
    echo ""
    echo "Saving commands for you..."
    
    # Create a file with setup commands
    cat > /tmp/betzone_setup_with_password.sh << 'SETUP_EOF'
#!/bin/bash
mysql -u root -p" {{ PASSWORD }}" << 'SQL_EOF'
CREATE DATABASE IF NOT EXISTS betzone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'betzone_password_2024';
GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';
FLUSH PRIVILEGES;
SHOW DATABASES;
SQL_EOF
SETUP_EOF
    
    # Replace password placeholder
    sed -i "s|{{ PASSWORD }}|$MYSQL_PASSWORD|g" /tmp/betzone_setup_with_password.sh
    chmod +x /tmp/betzone_setup_with_password.sh
    
    echo "Running database setup..."
    bash /tmp/betzone_setup_with_password.sh
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✓ Database setup successful!"
        echo ""
        echo "You can now run:"
        echo "  make dev"
    fi
    
else
    echo "✗ Password didn't work"
    echo ""
    echo "Not working? Try these options:"
    echo ""
    echo "OPTION 2: Reset MySQL root password"
    echo "   See: MYSQL_PASSWORD_HELP.md (Option B)"
    echo ""
    echo "OPTION 3: Manual setup"
    echo "   Run: mysql -u root -p"
    echo "   Then paste: MYSQL_SETUP_MANUAL.sql"
fi
