all:
	go build && ./skit
	#docker kill skit; docker rm skit
	#docker build -t skit .
	#docker run --name=skit -d -p 8080:8080 skit
