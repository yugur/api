# Yugur RESTful API
Source for the Yugur RESTful API.

## About
The API currently supports typical CRUD operations and basic auth, however auth is not fully integrated with the other handlers at this time (e.g. when you make a call to delete a dictionary entry, the API will allow it regardless of whether or not it received a valid session token.) You may also notice that there are two DB init scripts: `scripts/demo.sql` and `scripts/old_demo.sql`. The new script represents the DB structure that we are moving towards, while the old script still works out-of-the-box with the current source. For the time being you should use the old script.

### Handlers
- PATH w.r.t. / (METHODS...): Description
- status (GET): will return OK if API is running correctly, error otherwise.
- register (GET, POST): GET can be ignored if you have a frontend. POST is used to create a new user by providing the formvalues `username` and `password`
- login (GET, POST): as above but for auth. You will receive a cookie in response if login was successful.
- entry (GET, POST, PUT, DELETE): typical CRUD operations on individual dictionary entries. GET and DELETE both expect a formvalue `q` specifying the entry ID, while POST and PUT both expect a JSON object outlined below.
- fetch (GET): only relevant for testing/sanity checks. returns a JSON object of the entire dictionary.
- search-letter (GET): given a single letter in formvalue `q` this will return all entries beginning with that letter.
- search-tag (GET): given a tag ID, this will return all entries linked to that tag.

** Coming soon: **
- search (GET): given a query string (initially just a word) this will return all entries matching that word. This is the main way that a front end should obtain entry IDs that it can then use to direct the user to a specific page.
- imfeelinglucky (GET): returns a random entry.

### Dictionary objects
A dictionary entry is currently represented in JSON like so:
```
{
	"id":         "1",
	"headword":   "fire",
	"wordtype":   "1",
	"definition": "Burning fuel or other material: a cooking fire; a forest fire.",
	"hw_lang":    "1",
	"def_lang":   "1"
}
```
* `id` is unique to every entry.
* `headword` is what a user expects to see.
* `wordtype` is an ID corresponding to a row in the wordtypes table.
* `hw_lang` and `def_lang` are IDs corresponding to rows in the languages table. These are used to define the language of the headword and definition respectively. These may be changed soon to be unique locale codes like `en-AU` instead.

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
# \i scripts/old_demo.sql
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
The Yugur.io API is being developed by:

- Nicholas Brown (@nicholasbrown)
- Alexander Lawrence (@alexlawrence9)
