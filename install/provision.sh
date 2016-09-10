#!/bin/bash
set -e

#add go&java repo
add-apt-repository ppa:ubuntu-lxc/lxd-stable -y
add-apt-repository ppa:webupd8team/java -y
apt-get update

#install git
apt-get install git -y

#clone app
git clone https://github.com/pppepito86/smg.git /app/judge

#install mysql driver
git clone https://github.com/go-sql-driver/mysql.git /app/judge/src/github.com/go-sql-driver/mysql

#install go
apt-get install golang -y

#set GOPATH
echo "GOPATH=/app/judge" >> /etc/environment
echo "LC_ALL=en_US.UTF-8" >> /etc/environment
echo "LC_CTYPE=en_US.UTF-8" >> /etc/environment
echo "LANG=en_US.UTF-8" >> /etc/environment
echo "LANGUAGE=en_US.UTF-8" >> /etc/environment
source /etc/environment

#install mysql
echo "mysql-server-5.6 mysql-server/root_password password password" | debconf-set-selections
echo "mysql-server-5.6 mysql-server/root_password_again password password" | debconf-set-selections
apt-get install mysql-server-5.6 -y

#import database smg
mysql -u root -ppassword smg < smg.sql

#install docker
wget -qO- https://get.docker.com/ | sh
docker pull pppepito86/judgebox

#install unzip
apt-get install unzip -y

#install java
echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections
echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 seen true" | debconf-set-selections
apt-get install oracle-java8-installer -y