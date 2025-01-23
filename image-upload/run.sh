#!/bin/bash
# Build Docker image
echo -e "\033[1;38;5;155mBuilding Docker image... \033[0m"
docker image build -f Dockerfile -t forum .
echo -e "\033[1;38;5;155mList All Docker images... \033[0m"
docker images
# Run Docker container
echo -e "\033[1;38;5;155mRunning Docker container with Port 8080:8080 ... \033[0m"
docker container run -p 8404:8404 --detach --name forum forum
# Show running containers
echo -e "\033[1;38;5;155mRunning Containers: \033[0m"
docker ps
# Stop and remove container
docker image prune -a
# show files inside container
echo -e "\033[1;38;5;155mGo inside myContainers, list and exit: \033[0m"
docker exec -it forum /bin/bash -c "ls -l && exit"