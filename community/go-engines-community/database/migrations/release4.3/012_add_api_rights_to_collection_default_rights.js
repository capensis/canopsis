db.default_rights.insertMany([
    {
        _id: "api_app_info_read",
        loader_id: "api_app_info_read",
        crecord_name: "api_app_info_read",
        crecord_type: "action",
        desc: "read app info",
    },
    {
        _id: "api_user_interface_update",
        loader_id: "api_user_interface_update",
        crecord_name: "api_user_interface_update",
        crecord_type: "action",
        desc: "update user interface",
    },
    {
        _id: "api_user_interface_delete",
        loader_id: "api_user_interface_delete",
        crecord_name: "api_user_interface_delete",
        crecord_type: "action",
        desc: "delete user interface",
    },
    {
        _id: "api_broadcast_message",
        loader_id: "api_broadcast_message",
        crecord_name: "api_broadcast_message",
        crecord_type: "action",
        desc: "Broadcast Message",
        type: "CRUD",
    },
    {
        _id: "api_associative_table",
        loader_id: "api_associative_table",
        crecord_name: "api_associative_table",
        crecord_type: "action",
        desc: "Associative table",
        type: "CRUD",
    }
]);
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_app_info_read",
            "api_user_interface_update",
            "api_user_interface_delete"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_app_info_read"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.updateMany(
        {
            crecord_name: { "$in": ["Manager", "Support", "Visualisation", "Supervision"] },
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_associative_table"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.updateMany(
        {
            crecord_name: { "$in": ["Manager", "Support", "Visualisation", "Supervision"] },
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 4,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_broadcast_message",
            "api_associative_table"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 15,
                    crecord_type: "right",
                },
            },
        }
    )
});

var ui = db.configuration.find({'_id' : 'user_interface'})[0];
var d = {};
if (ui.max_matched_items === undefined) {
    d['max_matched_items'] = 10000
}
if (ui.check_count_request_timeout === undefined) {
    d['check_count_request_timeout'] = 30
}
if (Object.keys(d).length > 0) {
    db.configuration.update({'_id' : 'user_interface'},{$set : d})
}
