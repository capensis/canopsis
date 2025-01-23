db.external_data_table.createIndex({type: 1, name: 1}, {name: "type_1_name_1", unique: true});

if (!db.permission.findOne({_id: "api_external_data_table"})) {
    db.permission.insertOne({
        _id: "api_external_data_table",
        name: "api_external_data_table",
        description: "External data",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_external_data_table": 15
        }
    });
}
