#!/bin/bash

SERVER_URL="http://localhost:3000"

test_server_running() {
    # testing if the server is running...
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/")
    
    if [ "$RESPONSE" -eq 200 ]; then
        echo "✅ server is running (status: 200 ok)";
    else
        echo "❌ server is not running (received status: $RESPONSE)";
        exit 1;
    fi
}

test_submit_email() {
    # testing email submission...
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -d "email=test@example.com" "$SERVER_URL/submit")
    
    if [ "$RESPONSE" -eq 200 ]; then
        echo "✅ email submission successful (status: 200 ok)";
    else
        echo "❌ email submission failed (received status: $RESPONSE)";
        exit 1;
    fi
}

test_db_entry() {
    # testing if the email was written to the sqlite database...
    DB_FILE="sqlite-server/emails.db"
    
    if [ -f "$DB_FILE" ]; then
        if ! command -v sqlite3 >/dev/null 2>&1; then
            echo "❌ sqlite3 is not installed. please install sqlite3 to run this test";
            exit 1;
        fi;
        
        result=$(sqlite3 "$DB_FILE" "select email from emails where email='test@example.com';")
        if [ "$result" == "test@example.com" ]; then
            echo "✅ email found in $DB_FILE";
        else
            echo "❌ email not found in $DB_FILE";
            exit 1;
        fi;
    else
        echo "❌ database file $DB_FILE not found!";
        exit 1;
    fi
}

# run the tests
test_server_running;
test_submit_email;
sleep 2;  # wait for the write to complete
test_db_entry;

echo "✅ all tests passed successfully!";
