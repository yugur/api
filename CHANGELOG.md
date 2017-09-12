# Yugur API Changelog

## Upcoming

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


### Still to come
* Utility functions to speed up the id <=> name process (where relevant)
* A search handler to accomodate the fact that headwords are no longer unique
	* This isn't expected to perform any fancy IR, it will just grab all entries that have that headword.
* Debugging improvements, including new functions for the util package.
* (Possibly) functions or a package to replace the clunky in-line SQL.

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