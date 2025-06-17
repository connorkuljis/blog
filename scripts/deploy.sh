#!/bin/bash
set -e  # Exit immediately if a command exits with non-zero status

echo "Syncing local assets to staging directory..."

rsync --delete -avz \
	public/ \
	prod@kuljis.xyz:/home/prod/www/blog/public/

