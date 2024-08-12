# Lusca

Feature flag and configurations management system. It's named by a tentacled creature.

## Requirement

### golang-migrate

Detail: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Download pre-built binary

- MacOS

    ```bash
    $ brew install golang-migrate
    ```

- Windows, using https://scoop.sh/

    ```shell
    $ scoop install migrate
    ```

- Linux *.deb package

    ```bash
    $ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
    $ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
    $ apt-get update
    $ apt-get install -y migrate
    ```