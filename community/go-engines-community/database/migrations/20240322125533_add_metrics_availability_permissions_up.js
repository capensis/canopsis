if (!db.permission.findOne({_id: "availability_interval"})) {
    db.permission.insertOne({
        _id: "availability_interval",
        name: "availability_interval",
        description: "Availability interval"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_interval": 1
        }
    });
}

if (!db.permission.findOne({_id: "availability_listFilters"})) {
    db.permission.insertOne({
        _id: "availability_listFilters",
        name: "availability_listFilters",
        description: "Availability list filters"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_listFilters": 1
        }
    });
}

if (!db.permission.findOne({_id: "availability_editFilter"})) {
    db.permission.insertOne({
        _id: "availability_editFilter",
        name: "availability_editFilter",
        description: "Availability edit filter"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_editFilter": 1
        }
    });
}

if (!db.permission.findOne({_id: "availability_addFilter"})) {
    db.permission.insertOne({
        _id: "availability_addFilter",
        name: "availability_addFilter",
        description: "Availability add filter"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_addFilter": 1
        }
    });
}

if (!db.permission.findOne({_id: "availability_userFilter"})) {
    db.permission.insertOne({
        _id: "availability_userFilter",
        name: "availability_userFilter",
        description: "Availability user filter"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_userFilter": 1
        }
    });
}

if (!db.permission.findOne({_id: "availability_exportAsCsv"})) {
    db.permission.insertOne({
        _id: "availability_exportAsCsv",
        name: "availability_exportAsCsv",
        description: "Availability export as csv"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.availability_exportAsCsv": 1
        }
    });
}
