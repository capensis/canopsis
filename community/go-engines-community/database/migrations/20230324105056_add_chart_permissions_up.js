if (!db.default_rights.findOne({_id: "barchart_interval"})) {
    db.default_rights.insertOne({
        _id: "barchart_interval",
        crecord_name: "barchart_interval",
        crecord_type: "action",
        description: "Barchart interval",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_interval": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "barchart_sampling"})) {
    db.default_rights.insertOne({
        _id: "barchart_sampling",
        crecord_name: "barchart_sampling",
        crecord_type: "action",
        description: "Barchart sampling",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_sampling": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "barchart_listFilters"})) {
    db.default_rights.insertOne({
        _id: "barchart_listFilters",
        crecord_name: "barchart_listFilters",
        crecord_type: "action",
        description: "Barchart list filters",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_listFilters": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "barchart_editFilter"})) {
    db.default_rights.insertOne({
        _id: "barchart_editFilter",
        crecord_name: "barchart_editFilter",
        crecord_type: "action",
        description: "Barchart edit filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_editFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "barchart_addFilter"})) {
    db.default_rights.insertOne({
        _id: "barchart_addFilter",
        crecord_name: "barchart_addFilter",
        crecord_type: "action",
        description: "Barchart add filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_addFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "barchart_userFilter"})) {
    db.default_rights.insertOne({
        _id: "barchart_userFilter",
        crecord_name: "barchart_userFilter",
        crecord_type: "action",
        description: "Barchart user filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.barchart_userFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_interval"})) {
    db.default_rights.insertOne({
        _id: "linechart_interval",
        crecord_name: "linechart_interval",
        crecord_type: "action",
        description: "Linechart interval",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_interval": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_sampling"})) {
    db.default_rights.insertOne({
        _id: "linechart_sampling",
        crecord_name: "linechart_sampling",
        crecord_type: "action",
        description: "Linechart sampling",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_sampling": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_listFilters"})) {
    db.default_rights.insertOne({
        _id: "linechart_listFilters",
        crecord_name: "linechart_listFilters",
        crecord_type: "action",
        description: "Linechart list filters",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_listFilters": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_editFilter"})) {
    db.default_rights.insertOne({
        _id: "linechart_editFilter",
        crecord_name: "linechart_editFilter",
        crecord_type: "action",
        description: "Linechart edit filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_editFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_addFilter"})) {
    db.default_rights.insertOne({
        _id: "linechart_addFilter",
        crecord_name: "linechart_addFilter",
        crecord_type: "action",
        description: "Linechart add filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_addFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "linechart_userFilter"})) {
    db.default_rights.insertOne({
        _id: "linechart_userFilter",
        crecord_name: "linechart_userFilter",
        crecord_type: "action",
        description: "Linechart user filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.linechart_userFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_interval"})) {
    db.default_rights.insertOne({
        _id: "piechart_interval",
        crecord_name: "piechart_interval",
        crecord_type: "action",
        description: "Piechart interval",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_interval": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_sampling"})) {
    db.default_rights.insertOne({
        _id: "piechart_sampling",
        crecord_name: "piechart_sampling",
        crecord_type: "action",
        description: "Piechart sampling",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_sampling": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_listFilters"})) {
    db.default_rights.insertOne({
        _id: "piechart_listFilters",
        crecord_name: "piechart_listFilters",
        crecord_type: "action",
        description: "Piechart list filters",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_listFilters": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_editFilter"})) {
    db.default_rights.insertOne({
        _id: "piechart_editFilter",
        crecord_name: "piechart_editFilter",
        crecord_type: "action",
        description: "Piechart edit filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_editFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_addFilter"})) {
    db.default_rights.insertOne({
        _id: "piechart_addFilter",
        crecord_name: "piechart_addFilter",
        crecord_type: "action",
        description: "Piechart add filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_addFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "piechart_userFilter"})) {
    db.default_rights.insertOne({
        _id: "piechart_userFilter",
        crecord_name: "piechart_userFilter",
        crecord_type: "action",
        description: "Piechart user filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.piechart_userFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_interval"})) {
    db.default_rights.insertOne({
        _id: "numbers_interval",
        crecord_name: "numbers_interval",
        crecord_type: "action",
        description: "Numbers interval",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_interval": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_sampling"})) {
    db.default_rights.insertOne({
        _id: "numbers_sampling",
        crecord_name: "numbers_sampling",
        crecord_type: "action",
        description: "Numbers sampling",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_sampling": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_listFilters"})) {
    db.default_rights.insertOne({
        _id: "numbers_listFilters",
        crecord_name: "numbers_listFilters",
        crecord_type: "action",
        description: "Numbers list filters",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_listFilters": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_editFilter"})) {
    db.default_rights.insertOne({
        _id: "numbers_editFilter",
        crecord_name: "numbers_editFilter",
        crecord_type: "action",
        description: "Numbers edit filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_editFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_addFilter"})) {
    db.default_rights.insertOne({
        _id: "numbers_addFilter",
        crecord_name: "numbers_addFilter",
        crecord_type: "action",
        description: "Numbers add filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_addFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "numbers_userFilter"})) {
    db.default_rights.insertOne({
        _id: "numbers_userFilter",
        crecord_name: "numbers_userFilter",
        crecord_type: "action",
        description: "Numbers user filter",
        type: "RW"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.numbers_userFilter": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
