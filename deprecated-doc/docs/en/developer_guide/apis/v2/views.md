# Views API

## Views

### List views

```
GET /api/v2/views
GET /api/v2/views?name=...
GET /api/v2/views?title=...
```

Return all the views, optionally filtered by name or title.

```
{
    "groups": {
        "group1": {
            "name": "...",
            "views": [
                // ...
            ]
        },
        // ...
    }
}
```

### Create a view

```
POST /api/v2/views
{
    "group_id": "<group_id>",
    "type": "...",
    "name": "...",
    "title": "..."
    // ...
}
```

Create a view and return its id (which is generated automatically).

### Get a view

```
GET /api/v2/views/<view_id>
```

### Edit a view

```
PUT /api/v2/views/<view_id>
{
    "group_id": "<group_id>",
    "type": "...",
    "name": "...",
    "title": "..."
    // ...
}
```

### Remove a view

```
DELETE /api/v2/views/<view_id>
```


## Groups

### List groups

```
GET /api/v2/views/groups
GET /api/v2/views/groups?name=...
```

Return the list of all the groups, optionally filtered by name.

### Create a group

```
POST /api/v2/views/groups
{
    "name": "..."
}
```

Create a group and return its id (which is generated automatically).

### List the views of a group

```
GET /api/v2/views/groups/<group_id>
GET /api/v2/views/groups/<group_id>?name=...
GET /api/v2/views/groups/<group_id>?title=...
```

Return the list of the views of a group, optionally filtered by name or title.

```
{
    "_id": "...",
    "name": "...",
    "views": [
        ...
    ]
}
```

### Edit a group

```
PUT /api/v2/views/groups/<group_id>
{
    "name": "..."
}
```

### Remove a group

```
DELETE /api/v2/views/groups/<group_id>
```
