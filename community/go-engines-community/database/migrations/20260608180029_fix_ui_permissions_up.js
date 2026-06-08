db.permission.updateOne({_id: "listalarm_pbehavior"}, {
    $set: {
        api_pbehaviortype: 4,
        api_pbehaviorreason: 4,
        api_pbehaviorexception: 4,
    }
});

db.role.updateMany({"permissions.listalarm_pbehavior": 1}, {
    $set: {
        "permissions.api_pbehaviortype": 4,
        "permissions.api_pbehaviorreason": 4,
        "permissions.api_pbehaviorexception": 4,
    }
});

db.permission.updateOne({_id: "listalarm_fastPbehavior"}, {
    $set: {
        api_pbehavior: 9,
        api_pbehaviortype: 4,
        api_pbehaviorreason: 4,
    }
});

db.role.updateMany({"permissions.listalarm_fastPbehavior": 1}, {
    $set: {
        "permissions.api_pbehavior": 9,
        "permissions.api_pbehaviortype": 4,
        "permissions.api_pbehaviorreason": 4,
    }
});
