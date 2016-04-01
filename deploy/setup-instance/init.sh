#!/bin/bash
#
# Set up environment for mobile-push microservice

#
# ARGV
#
if [[ -z "$1" ]]; then
    echo "usage: $0 <region>"
    exit 1
fi

region=$1

#
# utilities
#
err() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%Sz')]: $@" >&2
}

#
# apt
#
apt-get update
apt-get install -y supervisor git make build-essential python-dev ntp

#
# install AWS CodeDeploy agent
#
apt-get install -y python-pip ruby2.0
pip install awscli
cd /home/ubuntu
wget https://aws-codedeploy-$region.s3.amazonaws.com/latest/install
chmod +x ./install
./install auto
rm ./install

#
# set up code directory structure
#
mkdir -p /srv/news-feed-writer/release
mkdir -p /srv/news-feed-writer/share
chown -R ubuntu:ubuntu /srv/news-feed-writer

mkdir -p /srv/news-feed-writer-staging/release
mkdir -p /srv/news-feed-writer-staging/share
chown -R ubuntu:ubuntu /srv/news-feed-writer-staging

#
# supervisor
#
cp ./news-feed-writer.conf /etc/supervisor/conf.d/
cp ./news-feed-writer-staging.conf /etc/supervisor/conf.d/
