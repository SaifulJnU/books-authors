# books-authors
This is ```REST-APi``` Project. In this project I implemented some request for Book and Athor.

Here is the Postman Documentaion: 
```https://documenter.getpostman.com/view/28855987/2s9YC7Sr2q```

---

To Run the Project in your local PC:
```docker build -t books-authors . ```
```sudo docker-compose up -d ```

---

Now to see the logs run the following command from the project root directory:
```docker logs books-authors_api_1```

---
To check the Prometheus type in browser:
```localhost:9090```
To see the number of successful login, place the following text(prometheus metric name) in the prometheus:
```myapp_successful_logins_total```

---

## More details are coming soon 
//mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0
//to run prometheus yml: saiful@saiful-Inspiron-3542:~/Downloads/prometheus-2.47.0.linux-amd64$ ./prometheus --config.file=prometheus.yml
//after that in browser use: localhost:9090
