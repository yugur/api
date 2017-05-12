# yugur API
API for the yugur.io crowd-sourced dictionary.

# Installation
Installation assumes a clean Ubuntu/Debian installation that already has Go 1.8 installed.

1. Install and setup psql
```
$ sudo apt-get install postgresql postgresql-contrib
$ sudo -i -u postgres
$ createdb yugur
$ psql
# ALTER USER postgres WITH PASSWORD 'postgres';
```

2. Initialize database

You may need to create a role for your account before running any scripts.
```
# \i demo.sql
```

3. Get the API

Using go get:
```
$ go get github.com/yugur/api
```

Using git:
```
$ cd $GOPATH/src/github.com/yugur
$ git clone git@github.com:yugur/api.git
```

4. Install the API
```
$ go install github.com/yugur/api
```

5. Run
```
$ api
```

Go to `localhost:8080` to make sure it works.

# Contribution
The Yugur.io API is being actively contributed to by:

- Nicholas Brown (@nicholasbrown)
- Alexander Lawrence (@alexlawrence9)