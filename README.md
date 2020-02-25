Start Services
- ``docker-compose up -d``
- ``swag init``
- ``go build main.go`` 
- ``go run main.go`` or ``ENVIRONMENT=local watcher``

Example: <br/>
`swag init && go run main.go`

Swagger:
https://gh-dark.rauchg.now.sh/swaggo/swag

##Database

- REFRESH MATERIALIZED VIEW {{name}}

- http://localhost:3003
- http://localhost:9200
- https://github.com/struCoder/pmgo
- http://localhost:8000/graphql


Crop images
- https://github.com/h2non/bimg

Install dependencies
- go get -u ./...



#### PRODUCTION
- sudo git pull
- pmgo kill
- pmgo start ./ golang --args="development"
