# Yugur RESTful API
Source for the Yugur RESTful API.

## About
The Yugur RESTful API provides a simple back-end for a dictionary service. 

### Handlers
- /status (GET): will return OK if API is running correctly, error otherwise.
- /search (GET): generic searching endpoint. Given a query via formvalue `q` this will compile a unique set of relevant results and return them as a JSON object. 
- /register (GET, POST): GET can be ignored if you have a frontend. POST is used to create a new user by providing the formvalues `username` and `password`
- /login (GET, POST): as above but for auth. You will receive a cookie in response if login was successful.
- /entry (GET, POST, PUT, DELETE): typical CRUD operations on individual dictionary entries. GET and DELETE both expect a formvalue `q` specifying the entry ID, while POST and PUT both expect a JSON object outlined below.
- /fetch (GET): only relevant for testing/sanity checks. returns a JSON object of the entire dictionary.
- letter (GET): given a single letter in formvalue `q` this will return all entries beginning with that letter.
- tag (GET): given a tag ID, this will return all entries linked to that tag.

### Dictionary objects
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
The fields `wordtype`, `hw_lang` and `def_lang` are stored in the database as a unique `bigint`, but they are converted to/from human-readable names as they exit/enter the API.


## Installation
The following instructions will get you started with a very simple setup. You may wish to modify the config file in which case your commands will differ (e.g. psql username and password). Installation assumes a clean Ubuntu/Debian installation that already has [Go 1.8](https://golang.org/dl/). 

1. Install and setup psql
```
$ sudo apt-get install postgresql postgresql-contrib
$ sudo -i -u postgres
$ createdb yugur
$ psql
# ALTER USER postgres WITH PASSWORD 'postgres';
```

2. Initialize database

If you want to run the script from the context of your own account (not the postgres account) then you will first need to create a new role with the same name as your own account (also a good idea if you want to run this on any internet-facing device. alter `config.json` accordingly.)
```
yourname$ (optional) sudo -u postgres createuser --interactive
yourname$ psql -d yugur
# \i scripts/demo.sql
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
or if you haven't added your bin directory to your `PATH`:
```
$ $GOPATH/bin/api
```

Go to `localhost:8080/status` to make sure it works

# Contribution
Development is led by [@nicholasbrown](https://github.com/nicholasbrown) with contributions from @alexlawrence9
