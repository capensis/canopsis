if (!db.permission.findOne({_id: "api_pbehavior_patterns"})) {
    db.permission.insertOne({
        _id: "api_pbehavior_patterns",
        name: "api_pbehavior_patterns",
        description: "Check PBehavior patterns",
        groups: ["api", "api_planning"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_pbehavior_patterns": 1
        }
    });
}
