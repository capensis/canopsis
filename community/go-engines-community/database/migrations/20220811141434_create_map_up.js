db.createCollection("map");

if (!db.default_rights.findOne({_id: "api_map"})) {
    db.default_rights.insertOne({
        _id: "api_map",
        crecord_name: "api_map",
        crecord_type: "action",
        desc: "Map",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_map": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "models_map"})) {
    db.default_rights.insertOne({
        _id: "models_map",
        crecord_name: "models_map",
        crecord_type: "action",
        desc: "Map",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_map": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_alarmsList"})) {
    db.default_rights.insertOne({
        _id: "map_alarmsList",
        crecord_name: "map_alarmsList",
        crecord_type: "action",
        desc: "Map: Access to 'Alarms list' modal"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_alarmsList": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_listFilters"})) {
    db.default_rights.insertOne({
        _id: "map_listFilters",
        crecord_name: "map_listFilters",
        crecord_type: "action",
        desc: "Rights on map: List filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_listFilters": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_editFilter"})) {
    db.default_rights.insertOne({
        _id: "map_editFilter",
        crecord_name: "map_editFilter",
        crecord_type: "action",
        desc: "Rights on map: Edit filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_editFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_addFilter"})) {
    db.default_rights.insertOne({
        _id: "map_addFilter",
        crecord_name: "map_addFilter",
        crecord_type: "action",
        desc: "Rights on map: Add filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_addFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_userFilter"})) {
    db.default_rights.insertOne({
        _id: "map_userFilter",
        crecord_name: "map_userFilter",
        crecord_type: "action",
        desc: "Rights on map: User filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_userFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "map_category"})) {
    db.default_rights.insertOne({
        _id: "map_category",
        crecord_name: "map_category",
        crecord_type: "action",
        desc: "Rights on map: Access to 'Category' action"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.map_category": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
