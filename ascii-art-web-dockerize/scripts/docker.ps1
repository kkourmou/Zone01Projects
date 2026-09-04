docker rm -f $(docker ps -aq)
docker rmi -f $(docker images -q)
docker image build -f Dockerfile -t ascii-art-web .
docker container run -p 8080:8080 -d --name ascii-art ascii-art-web