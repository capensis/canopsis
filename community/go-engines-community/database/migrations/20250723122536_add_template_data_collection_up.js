db.template_data.createIndex({name: 1}, {name: "name_1", unique: true});
db.template_test.createIndex({name: 1}, {name: "name_1", unique: true});
db.template_test.createIndex({type: 1, "rule._id": 1}, {name: "type_1_rule._id_1"});

if (!db.permission.findOne({_id: "api_template_data"})) {
    db.permission.insertOne({
        _id: "api_template_data",
        name: "api_template_data",
        type: "CRUD",
        description: "Template data",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_template_data": 15
        }
    });
}
