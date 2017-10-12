# Yugur RESTful API

The Yugur RESTful API aims to provide a simple, lightweight and multiplatform online dictionary API. It is designed to be completely independent of any frontend, however the [Yugur web app](https://github.com/yugur/yugurWebApp) is being developed to provide a default option for end-user access.

## Overview

### Dictionary entries

The API provides access to dictionary entries packaged as JSON objects.

```
{
	"id":         "1",
	"headword":   "fire",
	"wordtype":   "noun",
	"definition": "Burning fuel or other material: a cooking fire; a forest fire.",
	"hw_lang":    "en-AU",
	"def_lang":   "en-AU"
}
```

All communication with the API is done using either header form values or by including a JSON object like the one above in the body of the request.

### Endpoints

The main communication endpoints are:

* **status** - returns HTTP OK. In the future it will also return other useful status information in a JSON body.
* **search** - takes a query and returns a collection of (hopefully) relevant dictionary entries.
* **entry** - used to manipulate the dictionary entries by providing full Create, Read, Update, Delete access.
* **register** - used to register a new user with the API. Note that user accounts are extremely basic and currently have little function outside of authorisation.
* **login** - creates a new session and returns a cookie to the user if their login was successful.

There are more endpoints for manipulating components such as wordtypes and tags however these are still readily changing so they have not been included here for now.

## Getting Started

These instructions will get you a copy of the API up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the API on a live system.

### Prerequisites

These instructions assume that you are installing on a Linux distribution with access to `apt`. Currently there is no official guide for Windows or Mac OS. 

#### Go

The latest Go distributions can be downloaded [here](https://golang.org/dl/) and official instructions on how to install Go for your system are provided [here](https://golang.org/doc/install). These instructions assume that you have correctly set your Go PATH and have exported go/bin, for example:

```
$ nano ~/.bashrc
...
export PATH="$HOME/go/bin:$PATH"
...
```

#### Postgres

```
$ sudo apt-get install postgresql postgresql-contrib
```

### Installing

#### API configuration

It is recommended that you customise the database options contained in [config/config.json](config/config.json). Nevertheless the instructions will assume default settings which will otherwise set you up with a new Linux user `yugur` and corresponding psql role and database.

#### Compiling source

If you're familiar with Go this should be straightforward.

```
$ go install github.com/yugur/api
```

Note that the API expects to find [config/config.json](config/config.json) in its directory so you will need to move that with it. We will also be needing [scripts/](scripts) so copy that over too.

```
$ mkdir ~/yugur-api
$ cp $GOPATH/bin/api ~/yugur-api
$ mkdir ~/yugur-api/config
$ cp $GOPATH/src/github.com/yugur/api/config/config.json ~/yugur-api/config
```

#### Configuring the database

Create the database, add a new role for the API and set its password.

```
$ sudo -i -u postgres
$ createdb yugur
$ createuser yugur
$ psql
# ALTER USER yugur WITH ENCRYPTED PASSWORD 'yugur';
# \q
$ logout
```

#### Preparing UNIX user (optional)

This step can be skipped if you have setup a role and your settings to match your own UNIX user. However for deployment it is always a good idea to put services in their own non-superuser accounts. Create a new UNIX user for the API and move the files.

```
$ sudo adduser yugur
...
$ sudo cp -r ~/yugur-api /home/yugur
$ sudo chown -R yugur:yugur /home/yugur/yugur-api
```

#### Populating the database

The API does not currently execute any kind of "first start" code. It is assumed that you have manually setup everything beforehand. There is however a script that will initialise your database with the correct tables as well as some basic sample rows.

```
$ sudo -u yugur psql
# \i yugur-api/scripts/demo.sql
...
# \q
```

You can also find a shell script in the same folder that can be used to quickly add new dictionary entries. Fill a plain text file with an entry value one per line. That is, per dictionary entry you should have exactly five lines. For example, if you had two entries you might have a file `dict.txt` like so:

```
세상
noun
1.생명체가 살고 있는 지구 2.사람들이 생활하고 있는 사회 3.마음대로 활동할 수 있는 곳
ko-KR
ko-KR
세상
noun
1. world 2. planet 3. earth 4. era
ko-KR
en-AU
```

Then, run the script by specifying your data file and the IP address and port that your instance of the API is running on.

```
$ chmod +x populate.sh
$ ./populate.sh dict.txt localhost 8080
```

The script will then attempt to marshal your entries into JSON objects and send a POST request to the API's entry endpoint. The script will output how many requests it has processed as well as how they were represented in JSON, in case you are having issues with formatting your data file.

#### Run the API

The API does not require any special options, just run it from a command line.

```
$ sudo -u yugur yugur-api/api
...
The API is running at http://localhost:8080/
```

If nothing goes wrong you should see a series of positive sounding startup messages followed by the message above. You can check that everything is OK by going to `http://localhost:8080/status`. 

## Running tests

### Unit tests

The API includes unit tests for methods and functions where this is suitable. To run all included unit tests:

```
$ cd $GOPATH/src/github.com/yugur/api
$ go test ./...
```

Note that the crypto package tests can take a while due to hashing time requirements.

## Deployment

The API doesn't currently have built-in support to run as a service nor do we provide official binaries at this stage. As such you will want to follow the instructions above and then run it inside a detachable screen for any long term usage.

```
$ sudo apt-get install screen
$ screen
$ path/to/api
...
```

You can detach from the screen with `CTRL+A-D`

## Contributing

TBD

## Versioning

TBD

## Authors

* **Nicholas Brown** - [@nicholasbrown](https://github.com/nicholasbrown)

See also the list of [contributors](CONTRIBUTORS.md) that have participated in the project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* The community of native Western Yugur speakers for inspiring this project.
* This great [README template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2) by [@PurpleBooth](https://github.com/PurpleBooth)
