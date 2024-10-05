#!/bin/bash

make build

REMOTE_HOST="otso@piikki.stadi.ninja"
PROJECT_DIR="$HOME/piikittaja"
BINARY_DIR="$PROJECT_DIR/bin"
TAR_FILE="deploy_files.tar.gz"

# Kill existing processes (ignore errors if they are not running)
ssh $REMOTE_HOST "pkill -f web_app || true"
ssh $REMOTE_HOST "pkill -f telegram_bot || true"

# Create a tarball of the bin directory and assets excluding .env
tar --exclude="./assets/.env" -czf $TAR_FILE ./bin ./assets

# Copy the tarball to the remote host
scp $TAR_FILE $REMOTE_HOST:"$PROJECT_DIR"/

# Connect to the remote host, extract the tarball, and remove it
ssh $REMOTE_HOST << EOF
    tar -xzf "$PROJECT_DIR/$TAR_FILE" -C "$PROJECT_DIR"
    rm "$PROJECT_DIR/$TAR_FILE"
EOF

# Remove the local tarball after successful copy
rm $TAR_FILE

# Ensure remote log directory exists
ssh $REMOTE_HOST "mkdir -p $PROJECT_DIR/logs"

# Start both telegram bot and web app in new background processes
ssh $REMOTE_HOST "nohup $BINARY_DIR/web_app > $PROJECT_DIR/logs/web_app.log 2>&1 &"
ssh $REMOTE_HOST "nohup $BINARY_DIR/telegram_bot > $PROJECT_DIR/logs/telegram_bot.log 2>&1 &"
