# Custom SQL DB API service

SQL micro service is a RESTful API service that provides the following functionality:


## SQL API

    /sql
    Run sql commands
## payloads 
### create
```shell
curl --request POST \
  --url http://127.0.0.1:8080/query \
  --header 'Content-Type: application/json' \
  --data '{
	"query": "create table user (id int, name String)"
}'
```
### insert 
```shell
curl --request POST \
  --url http://127.0.0.1:8080/query \
  --header 'Content-Type: application/json' \
  --data '{
	"query": "insert into user values (2, '\''John'\'')"
}'
```

### select
```shell
curl --request POST \
  --url http://127.0.0.1:8000/query \
  --header 'Content-Type: application/json' \
  --data '{
	"query": "select * from user"
}'
```

### drop table
```shell
curl --request POST \
  --url http://127.0.0.1:8080/query \
  --header 'Content-Type: application/json' \
  --data '{
	"query": "drop table user"
}'
```



## Technologies & Tools Used
[Docker](https://www.docker.com/)
[Docker Compose](https://docs.docker.com/compose/overview/)
[GoLang](https://golang.org/doc/default_packages/net.html)
[Mockery](https://github.com/vektra/mockery)
[Git](https://git-scm.com/)
[VSCode](https://code.visualstudio.com/)
[Redis](https://redis.io/)




## Run the service

There is a Dockerfile in the parent directory of the project.
To run the service:
- go to the project directory 
- rename .env.sample to .env
- run `docker-compose up`
- project will be up on port 8000

## Test the service
To test the service:
run from root directory `go test ./... -v`
