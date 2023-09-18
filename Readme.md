# books-authors

This is a REST API project that implements various requests for Books and Authors.

[![Postman Documentation](https://img.shields.io/badge/Postman-Documentation-orange)](https://documenter.getpostman.com/view/28855987/2s9YC7Sr2q)

## Getting Started

To run the project on your local PC, follow these steps:

1. Build the Docker image: <br>
```bash
docker build -t books-authors .   
```
2. Start the containers in detached mode using Docker Compose: <br>
```bash
sudo docker-compose up -d
```

## Viewing Logs

To view the logs of the running containers, execute the following command from the project's root directory: <br>

```bash
docker logs books-authors_api_1
```

## Prometheus Metrics
To check Prometheus metrics, open your browser and navigate to: <br>
```localhost:9090```  <br>

## To monitor the number of successful logins, use the following Prometheus metric name: <br>

```bash
myapp_successful_logins_total
```

## Additional Details
More details about the project are coming soon.

