# Upgrading to Canopsis 3.2

## New rights

3 new rights have been introduced:
- `listalarm_editFilter`
- `listalarm_listFilters`
- `listalarm_addFilter`

Insert them through MongoDB:
```
mongodb:PRIMARY> db.default_rights.insert([
{"_id" : "listalarm_listFilters",     "loader_id" : "listalarm_listFilters",     "crecord_name" : "listalarm_listFilters",     "crecord_type" : "action",     "desc": "Rights on listalarm: list filters"},
{"_id" : "listalarm_editFilter",     "loader_id" : "listalarm_editFilter",     "crecord_name" : "listalarm_editFilter",     "crecord_type" : "action",     "desc": "Rights on listalarm: edit filters"},
{"_id" : "listalarm_addFilter",     "loader_id" : "listalarm_addFilter",     "crecord_name" : "listalarm_addFilter",     "crecord_type" : "action",     "desc": "Rights on listalarm: add filter"}
])
```

Then, go to the Canopsis UI as an admin user, and allow your users to use these new roles, if necessary.

*WARNING:* some browser cache might prevent you from using these new roles. Empty your cache, or reconnect into your user account.
