# Reverse-Proxy Cache Redis

Implementation of a reverse proxy that caches the database in Redis memory. Background update when the cache expires and scalable is some of the differentials.

### Prerequisites

* Install [Golang](https://golang.org/dl/)
* Install go-redis dependency: ```go get -u github.com/go-redis/redis```

### Configuration

All configuration is made by environment variables:
* **GATEWAY_HOST**: Host gataway. Ex: http://localhost:1333
* **PORT**: Port to exposes reverse-proxy. Ex: 80
* **REDIS_HOST**: Host of redis. Ex: localhost
* **REDIS_PORT**: Port of redis. Ex: 6379
* **CACHE_EXP**: TTL of cache expiration in seconds. Ex: 60

## Deployment

Run go build using ```go build``` and execute the binary file.

## Built With

* [Golang](https://golang.org/) - The main language used
* [go-redis](https://github.com/go-redis/redis) - Dependency of redis for golang

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/leoribeirowebmaster/Reverse-Proxy-Cache-Redis/tags). 

## Authors

* **Leonardo Ribeiro** - *Initial work* - [leoribeirowebmaster](https://github.com/leoribeirowebmaster)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details