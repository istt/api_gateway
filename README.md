# Development

### Bring up MongoDB via Docker

```bash
$  docker-compose -f deployments/docker/mongodb.yml up -d
Creating docker_jhipster-mongodb_1 ... done
```

Login to mongodb shell


```bash
$ docker exec -it docker_jhipster-mongodb_1 mongo
```

