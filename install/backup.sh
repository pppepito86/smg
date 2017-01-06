#!/bin/bash

date
mkdir -p /app/backup
mysqldump -u root -ppassword smg > /app/backup/smg.sql
cp -R /app/judge/src/workdir /app/backup
tar -zcf /app/backup.tar.gz /app/backup

aws glacier upload-archive --vault-name judge --account-id - --body /app/backup.tar.gz

rm -rf /app/backup
