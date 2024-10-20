#!/bin/bash

ENV_FILE="./assets/.env"

# Create or overwrite the .env file
touch $ENV_FILE

# Write environment variables to the .env file in one block
{
  echo "DATABASE_ADMIN=$DATABASE_ADMIN"
  echo "DATABASE_APP=$DATABASE_APP"
  echo "TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN"
  echo "GOOGLE_KEY=$GOOGLE_KEY"
  echo "GOOGLE_SECRET=$GOOGLE_SECRET"
  echo "COOKIE_STORE_SECRET=$COOKIE_STORE_SECRET"
  echo "HOST_URL=$HOST_URL"
  echo "PORT=$PORT"
} >> $ENV_FILE
