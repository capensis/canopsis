db.template_data.dropIndex("name_1");
db.template_test.dropIndex("name_1");
db.template_test.dropIndex("type_1_rule._id_1");

db.permission.deleteMany({
    _id: {
        $in: [
            "api_template_data",
            "models_templateTesting",
        ]
    }
});

db.role.updateMany({}, {
    $unset: {
        "permissions.api_template_data": "",
        "permissions.models_templateTesting": "",
    }
});
