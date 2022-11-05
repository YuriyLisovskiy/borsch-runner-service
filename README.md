## Trainee Golang Test Task
[![Test](https://github.com/YuriyLisovskiy/borsch-runner-service/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/YuriyLisovskiy/borsch-runner-service/actions/workflows/ci.yml)

### Deployment with Docker
Deploy the app:
```shell
docker compose up
```
or:
```shell
docker compose run app --rm --detach
```

The reference for `docker compose run` command options can be found
[here](https://docs.docker.com/engine/reference/commandline/compose_run/#options).

To stop containers, press `Ctrl+C` or run the command below if you started them in detached mode:
```shell
docker compose down
```

Perform cleaning up, if required:
```shell
docker compose rm --force --stop
```

To learn more about `docker compose rm` options, read the
[reference](https://docs.docker.com/engine/reference/commandline/compose_rm/#options).
