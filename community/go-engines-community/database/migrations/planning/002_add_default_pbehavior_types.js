function genID() {
    return UUID().toString().split('"')[1]
}

db.pbehavior_type.insertMany([
    {
        "_id": genID(),
        "name": "Default inactive",
        "description": "Default inactive",
        "icon_name": "brightness_3",
        "priority": 0,
        "type": "inactive"
    },
    {
        "_id": genID(),
        "name": "Default active",
        "description": "Default active",
        "icon_name": "",
        "priority": 1,
        "type": "active"
    },
    {
        "_id": genID(),
        "name": "Default maintenance",
        "description": "Default maintenance",
        "icon_name": "build",
        "priority": 2,
        "type": "maintenance"
    },
    {
        "_id": genID(),
        "name": "Default pause",
        "description": "Default pause",
        "icon_name": "pause",
        "priority": 3,
        "type": "pause"
    }
]);