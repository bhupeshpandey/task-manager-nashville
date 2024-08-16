# task-manager-nashville
This service handles the rest api's by creating a gropc client for the gallatin grpc server.

For this service just run the main.go locally.

go run /cmd/main.go

Also, this is a combination service. The Nashville, gallatin and ashland. For running all three together<
please run the docker-compose from the gallatin which contains rabbitmq, postgres, and redis cache.
Comment out the gallatin section in docker-compose as that is having some networking issues.

Once the services are up and running on docker,
Go ahead and execute /cmd/main.go on the nashville, gallatin and the ashland.
Once all three services are up, use the postman collection in the
to call the api's. 