# arvanDiscount

This project is my code challenge for ArvanCloud company interview.<br>
It uses redis pub/sub as a queue for communicate with wallet to apply credit.<br>
To run the project make sure you have installed docker and docker-compose and it's running.<br>
Use following command to check docker installed and it's running:
``` shell script
$docker -v
```
You should see something like:
``` shell script
 Docker version 19.03.13, build 4484c46d9d
```

After that check that you have installed docker-compose as well with following command:
```shell script
$docker-compose -v
```
You should see something like:
``` shell script
docker-compose version 1.27.4, build 40524192
```

Make a copy from `env.yaml.default` and rename it to `env.yml`:
```shell script
$cp env.yaml.default env.yaml
```

After that you have to use following command to run project:
```shell script
$docker-compose up -d
```

Wait till application starts after that use following command to see forwarded ports:
```shell script
$docker-compose ps
```
