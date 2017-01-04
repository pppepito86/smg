#!/bin/bash

# install mail utils
debconf-set-selections <<< "postfix postfix/mailname string pesho.org"
debconf-set-selections <<< "postfix postfix/main_mailer_type string 'Internet Site'"
apt-get install -y mailutils
sed -i 's/^myhostname.*$/myhostname = pesho.org/g' /etc/postfix/main.cf
sed -i 's/^myorigin.*$/myorigin = $myhostname/g' /etc/postfix/main.cf
sed -i 's/^inet_interfaces.*$/inet_interfaces = localhost/g' /etc/postfix/main.cf
# dkim tools - email dkim security
apt-get install -y opendkim opendkim-tools
cat <<EOF >> /etc/opendkim.conf
AutoRestart             Yes
AutoRestartRate         10/1h
UMask                   002
Syslog                  yes
SyslogSuccess           Yes
LogWhy                  Yes

Canonicalization        relaxed/simple

ExternalIgnoreList      refile:/etc/opendkim/TrustedHosts
InternalHosts           refile:/etc/opendkim/TrustedHosts
KeyTable                refile:/etc/opendkim/KeyTable
SigningTable            refile:/etc/opendkim/SigningTable

Mode                    sv
PidFile                 /var/run/opendkim/opendkim.pid
SignatureAlgorithm      rsa-sha256

UserID                  opendkim:opendkim

Socket                  inet:12301@localhost
EOF
cat <<EOF >> /etc/default/opendkim
SOCKET="inet:12301@localhost"
EOF
cat <<EOF >> /etc/postfix/main.cf
milter_protocol = 2
milter_default_action = accept
smtpd_milters = inet:localhost:12301
non_smtpd_milters = inet:localhost:12301
EOF
mkdir /etc/opendkim
mkdir /etc/opendkim/keys
cat <<EOF > /etc/opendkim/TrustedHosts
127.0.0.1
localhost
192.168.0.1/24

*.pesho.org
EOF
cat <<EOF > /etc/opendkim/KeyTable
default._domainkey.pesho.org pesho.org:default:/etc/opendkim/keys/pesho.org/default.private
EOF
cat <<EOF > /etc/opendkim/SigningTable
*@pesho.org default._domainkey.pesho.org
EOF
# generate keys
mkdir -p /etc/opendkim/keys/pesho.org
opendkim-genkey -s default -d pesho.org -D /etc/opendkim/keys/pesho.org/
chown opendkim:opendkim /etc/opendkim/keys/pesho.org/default.private
service postfix restart
service opendkim restart
# TLS Mail Encryption - Postfix configuration for MAY:(opportunistic)
cat <<EOF >> /etc/postfix/main.cf
smtp_tls_security_level = may
smtp_tls_policy_maps = hash:/etc/postfix/tls_policy
EOF
cat <<EOF > /etc/postfix/tls_policy
pesho.org may
.pesho.org may
EOF
postmap /etc/postfix/tls_policy
/etc/init.d/postfix reload
