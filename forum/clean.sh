#!/bin/bash
# Stop all running containers
echo -e "\033[1;38;5;155mStopping all Docker containers... \033[0m"
docker stop $(docker ps -q)
# Remove all stopped containers
echo -e "\033[1;38;5;155mRemoving all stopped Docker containers... \033[0m"
docker rm $(docker ps -a -q)
# Clean up dangling images
echo -e "\033[1;38;5;155mCleaning up Docker images... \033[0m"
docker image rm -f $(docker image ls -q)
echo -e "\033[1;38;5;155mClean-up completed. \033[0m"