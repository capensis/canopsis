db.role.updateMany({"$or":[
    {"permissions.models_exploitation_linkRule": {$ne: null}},
    {"permissions.models_exploitation_eventFilter": {$ne: null}},
]}, {
    $set: {
        "permissions.api_external_data_table": 4
    }
});
