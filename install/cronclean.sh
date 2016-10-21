#!/bin/bash

# Cron expression
cron="06 22 * * * /app/judge/install/clean.sh >> /app/clean.out 2>&1"

# Escape all the asterisks so we can grep for it
cron_escaped=$(echo "$cron" | sed s/\*/\\\\*/g)

crontab -l | grep "${cron_escaped}"
if [[ $? -eq 0 ]] ;
  then
    echo "Crontab already exists."
  else
    # Write out current crontab into temp file
    crontab -l > mycron
    # Append new cron into cron file
    echo "$cron" >> mycron
    # Install new cron file
    crontab mycron
    # Remove temp file
    rm mycron
fi
