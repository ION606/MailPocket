#!/bin/bash

SERVER_URL="http://localhost:3000"

test_server_running() {
    echo "Testing if the server is running..."
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/")
    
    if [ "$RESPONSE" -eq 200 ]; then
        echo "✅ Server is running (Status: 200 OK)"
    else
        echo "❌ Server is NOT running (Received status: $RESPONSE)"
        exit 1
    fi
}


test_submit_email() {
    echo "Testing email submission..."
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST -d "email=test@example.com" "$SERVER_URL/submit")
    
    if [ "$RESPONSE" -eq 200 ]; then
        echo "✅ Email submission successful (Status: 200 OK)"
    else
        echo "❌ Email submission failed (Received status: $RESPONSE)"
        exit 1
    fi
}


test_csv_written() {
    CSV_FILE="batched-server/emails.csv"
    
    if [ -f "$CSV_FILE" ]; then
        if grep -q "test@example.com" "$CSV_FILE"; then
            echo "✅ Email found in $CSV_FILE"
        else
            echo "❌ Email NOT found in $CSV_FILE"
            exit 1
        fi
    else
        echo "❌ CSV file not found!"
        exit 1
    fi
}

# run the tests
test_server_running
test_submit_email
sleep 2  # wait for the batched write to complete
test_csv_written

echo "✅ All tests passed successfully!"
