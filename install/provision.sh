#!/bin/bash

set -e -x

#add go&java repo
add-apt-repository ppa:ubuntu-lxc/lxd-stable -y
add-apt-repository ppa:webupd8team/java -y
apt-get update

#install git&stuff
apt-get install -y curl git gcc make python-dev vim-nox jq cgroup-lite silversearcher-ag

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

# Set up vim for golang development
git clone https://github.com/luan/vimfiles.git $HOME/.vim
curl vimfiles.luan.sh/install | bash

# set up bash-it
git clone --depth=1 https://github.com/Bash-it/bash-it.git $HOME/.bash_it
~/.bash_it/install.sh --silent
echo "echo -e -n \"\x1b[\x35 q\"" > ~/.bash_it/custom/cursor.sh

# Set up tmux
wget -O $HOME/.tmux.conf https://raw.githubusercontent.com/luan/dotfiles/master/tmux.conf

#install mysql
echo "mysql-server-5.6 mysql-server/root_password password password" | debconf-set-selections
echo "mysql-server-5.6 mysql-server/root_password_again password password" | debconf-set-selections
apt-get install mysql-server-5.7 -y

#import database smg
mysql -u root -ppassword < /app/judge/install/smg.sql

#install docker
wget -qO- https://get.docker.com/ | sh
docker pull pppepito86/judgebox

#install unzip
apt-get install unzip -y

#install java
echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections
echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 seen true" | debconf-set-selections
apt-get install oracle-java8-installer -y

#remove docker memory swap warning
#WARNING: Your kernel does not support swap limit capabilities, memory limited without swap.
#sed -i -e 's/GRUB_CMDLINE_LINUX=.*$/GRUB_CMDLINE_LINUX=\"cgroup_enable=memory swapaccount=1\"/' /etc/default/grub
#echo "GRUB_CMDLINE_LINUX=\"cgroup_enable=memory swapaccount=1\"" >> /etc/default/grub
#sh -c exec grub-mkconfig -o /boot/grub/grub.cfg "$@"

#create start service
cp /app/judge/install/judge /etc/init.d/judge

chmod 700 /etc/init.d/judge
update-rc.d judge defaults
update-rc.d judge enable

#mail
chmod +x /app/judge/install/mail.sh
/app/judge/install/mail.sh

#backup
chmod +x /app/judge/install/backup.sh
chmod +x /app/judge/install/cronbackup.sh
/app/judge/install/cronbackup.sh

#clean
chmod +x /app/judge/install/clean.sh
chmod +x /app/judge/install/cronclean.sh
/app/judge/install/cronclean.sh

mkdir -p /app/aws
curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o /app/aws/"awscli-bundle.zip"
unzip /app/aws/awscli-bundle.zip -d /app/aws
/app/aws/awscli-bundle/install -b /usr/bin/aws

#restart
sudo reboot

