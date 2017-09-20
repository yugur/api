# Yugur API Changelog

## 2017-09-20

### New features
* Generic search handler
	* Given a single token query q, this handler will compile a unique set of entries from all of the available searching methods.
	* No indexing, this won't scale but it will do for demonstration.
* Query functions
	* Most of the SQL functionality has been moved to a separate file.
	* Handlers make calls as appropriate to these functions.
	* This has cleaned up most of the handler code to be much more readable.
* Human-readable <-> database conversion
	* Entries are now processed on arrival/departure.
	* Basically a front-end no longer has to be aware of any concept of unique IDs (except for the entry ID)
	* Can be extended in the future to sanitise incoming requests.
* Error calls
	* Initial base for a generic error function.
	* Uses Go's ability to pass multi-paramater functions as the parameters for another function.
	* Subject to change.

### Changes
* Codebase has been updated to properly support the new database. All existing handlers should work, however many of their usages have been changed.
	* The entry handler now expects to receive an entry ID instead for all methods that previously relied on a headword.
	* Login handler still expects a username as before.
	* Register handler now expects an email address as well as a username/password.
	* Letter search handler still expects a letter as before.
	* Tag search handler now expects the tag ID instead of the name.
* Tag search handler has been fixed. The previous query didn't seem to be working for any of the demo scripts so it has been replaced with a new query that supports the entry_tags 1-to-many relationship.
* There is now an email field on the template register page.
* Sessions are now stored by ID instead of username.
* CORS is now enabled by default.

### Failed experiments
* An attempt to turn Go into Java by creating a special database package for all of the different queries. Seems possible but more complicated than it would be worth.

## 2017-08-28

### New features

* Added a **config package.**
	* Loads in config values from file.
	* Currently supports JSON.
	* The main package has been updated to take advantage of this.

### Changes
* Moved database and init to main.go
* Added some basic initialisation messages to aid troubleshooting.
* Updated the demo SQL script.
* Refactoring to move more towards idiomatic Go.

## 2017-08-08
* The changelog is born.