#!/usr/bin/env bash
echo "\n make build for target..."
GOOS="linux" GOARCH="amd64" go build -o ./bin/amd64 ./
echo "build complete! \n"
echo "sending file..."
sshpass -p "123123" scp -P9844 -r ./bin reza@188.40.164.251:/home/reza/
echo "file sent! \n"
echo "moving to /hellno directory"
sshpass -p "123123" ssh -t -p 9844 reza@188.40.164.251 "echo 123123 | sudo cp -r /home/reza/bin/* /hellno/reza/bin && rm -R /home/reza/*"
echo "build docker file..."
sshpass -p "123123" ssh -t -p 9844 reza@188.40.164.251 "cd /hellno/reza/ && sudo docker build -t reza ."
echo "docker build done!"
echo "create container..."
sshpass -p "123123" ssh -t -p 9844 reza@188.40.164.251 "cd /hellno/reza/ && docker-compose down && docker-compose up -d"
echo "container created!"
echo "\n complete!"