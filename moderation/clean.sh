#!/bin/bash
# Stop all running containers
echo -e "\033[1;38;5;155mStopping all Docker containers... \033[0m"
docker stop $(docker ps -q)
# Remove all stopped containers
echo -e "\033[1;38;5;155mRemoving all stopped Docker containers... \033[0m"
docker rm $(docker ps -a -q)
# Clean up dangling images
docker image prune -a -f
# Clean up unused volumes
echo -e "\033[1;38;5;155mCleaning up unused Docker volumes... \033[0m"
docker volume prune -f
echo -e "\033[1;38;5;155mClean-up completed. \033[0m"