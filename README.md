# The yuni language

The programming language that is creating in my live streaming( [プログラミング言語を作ろう](https://www.youtube.com/playlist?list=PL0GVAkk6KEQxuSyud9b81t7DKK83y2tXj) ).

## How to compile the compiler?

You can build the compiler in container by using [docker](https://www.docker.com/) and docker-compose.
After install `docker` and `docker-compose`, run following command on repository root:

```sh
docker compose run test
```

It will run compiler tests along with `test.sh`'s definition so you can see the result of tests.
If it is all SUCEEDED, you can dive into the container by using following command:

```sh
docker compose run sh
```

It takes TTY and so you can operate it with arbitrary unix commands.
The final binary is `./bin/lang`.
It takes STDIN and it will generate LLVM-IR to STDOUT so if you want to run a file then you can use `cat (anyfile) | ./bin/lang | lli` in the container.

## Structure of compiler environment

This slide page describes such info: https://docs.google.com/presentation/d/1GUWQv3kVH8Kv1apoQ1Gu01DG1EWngMzwSfiiCE46B3A/edit#slide=id.g112d216b7f3_0_44

## Examples

You can see examples on [./test/](./test) directory.
