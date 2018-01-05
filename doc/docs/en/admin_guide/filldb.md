# The filldb command

You can provision any collection of the Canopsis Database with this tool.

This command is useful to  clone a part of the Canopsis configuration from an evironment to another (ex: from pre-production to production)

## preparing json files to load in DB

to update a collection, `filldb` looks for a folder named `json_<collection_name> `  in
`~/opt/mongodb/load.d/ `

`filldb` will load all the json files from this folder in the collection defined by `json_<collection_name>`

### json files structure

each json file can contain an object that will create one entry in the collection, or an array that will create several entries in te collection.


> **NOTE** : you need to add 3 fields to each json object for filldb to work :

- `_id ` : an unique ID
- ` loader_id` : set it to the "id" value that can been seen on the canopsis UI
- `loader_no_update` : a boolean
	- set it to `true` to avoid updating an item if it already exists
	- set it to `false` to update the item even if it already exists in the DB


## usage

> **WARNING**: this command  **WILL DESTROY DATA**. Use it with caution.

`canopsis-filldb --init` :

- destroy all collections described in `~/etc/migration/purge.conf ` 

- loads all files in the ` ~/opt/mongodb/load.d/ `


`canopsis-filldb --update` :

- loads all files in the ` ~/opt/mongodb/load.d/ `
