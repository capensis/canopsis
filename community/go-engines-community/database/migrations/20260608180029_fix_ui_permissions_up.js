db.permission.updateOne({_id: "listalarm_pbehavior"}, {
    $set: {
        "api_permissions.api_pbehaviortype": 4,
        "api_permissions.api_pbehaviorreason": 4,
        "api_permissions.api_pbehaviorexception": 4,
    }
});

db.role.updateMany({"permissions.listalarm_pbehavior": 1}, {
    $bit: {
        "permissions.api_pbehaviortype": {or: 4},
        "permissions.api_pbehaviorreason": {or: 4},
        "permissions.api_pbehaviorexception": {or: 4},
    }
});

db.permission.updateOne({_id: "listalarm_fastPbehavior"}, {
    $set: {
        "api_permissions.api_pbehavior": 9,
        "api_permissions.api_pbehaviortype": 4,
        "api_permissions.api_pbehaviorreason": 4,
    }
});

db.role.updateMany({"permissions.listalarm_fastPbehavior": 1}, {
    $bit: {
        "permissions.api_pbehavior": {or: 9},
        "permissions.api_pbehaviortype": {or: 4},
        "permissions.api_pbehaviorreason": {or: 4},
    }
});
