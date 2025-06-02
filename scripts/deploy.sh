#!/bin/bash
set -e  # Exit immediately if a command exits with non-zero status

# Variables
REMOTE_HOST="${SSH_USER}@${SERVER_IP}"
STAGING_DIR="~/public/"
DOCUMENT_ROOT="/var/www/html/kuljis.xyz"
LOCAL_DIR="www_root/"

# Copy local assets to remote staging directory
echo "Syncing local assets to staging directory..."
rsync --delete -avz ${LOCAL_DIR} ${REMOTE_HOST}:${STAGING_DIR}

# Use a single SSH connection with sudo privileges
echo "Deploying from staging to document root..."

ssh -t ${REMOTE_HOST} \
	"sudo bash -c 'rm -rf ${DOCUMENT_ROOT}/* && cp -r ${STAGING_DIR}* ${DOCUMENT_ROOT}/'"

echo "Deployment completed successfully!"
