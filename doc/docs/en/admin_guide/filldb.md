# The filldb command

You can provision any collection of the Canopsis Database with this tool.

This command is useful to  clone a part of the Canopsis configuration from an evironment to another (ex: from pre-production to production)

## Preparing json files to load in DB

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

## Usage

> **WARNING**: this command  **WILL DESTROY DATA**. Use it with caution.

`canopsis-filldb --init` :

- destroy all collections described in `~/etc/migration/purge.conf ` 

- loads all files in the ` ~/opt/mongodb/load.d/ `


`canopsis-filldb --update` :

- loads all files in the ` ~/opt/mongodb/load.d/ `

## json loader module

The file module `jsonloader.py` enables to tell canopsis to load
custom files into the database depending in upsert mode following below
considerations:

The json loader script searches each folder located in the
`/opt/canopsis/opt/mongodb/load.d` folder that are prefixed by
**[json]()**. Theses folders name have to be followed by a collection
name that will tell where to upsert json documents. If a json document
have to be upsert in the object collection, the json folder have to be
called **json_object**.

Json folders to contain **<filename>.json** files that contains either
a document or a list of document. These document must be identified by a
special key that allow upsert. This key is named `loader_id` and must be
uniq over the collection folder documents. Below a sample of a json
document being inserted in the object collection because the file path
is `/opt/canopsis/opt/mongodb/load.d/json_object/mydocument.json`

in a single document:

```javascript
{
   "loader_id":"000",
   "document_key_1": "document_value_1"
}
```

or in a document list:

```javascript
[
   {
      "loader_id":"000",
      "document_key_2": "document_value_2"
   },
   {
      "loader_id":"000",
      "document_key_3": "document_value_3"
   }
]
```

it is possible to prevent a document for being updated by the json
loader by adding the `loader_no_update` key equals to **true** in the
json document.

### json loader hook

When the json loader is about to upsert a json document, some processing
is called. The feature this brings is for example to replace a macro
with a specific computed value.

> -   Macro `[[HOSTNAME]]` will tell the server where filldb is ran to
>     replace strings in records with the hostname value in place of the
>     macro. for example, when the following record is processed :

```javascript
{
   "loader_id":"000",
   "document_key_1": "document_value_1"
   "my_key" : "canopsis runs on [[HOSTNAME]]"
   "my_list" : ["canopsis runs on [[HOSTNAME]]"]
}
```

Will upsert in database the following document in case `myhostname` is
the server hostname:

```javascript
{
   "loader_id":"000",
   "document_key_1": "document_value_1"
   "my_key" : "canopsis runs on myhostname"
   "my_list" : ["canopsis runs on myhostname"]
}
```

It is possible in json documents that the jsonloader will proceed to use
a special macro that will replace record string in
