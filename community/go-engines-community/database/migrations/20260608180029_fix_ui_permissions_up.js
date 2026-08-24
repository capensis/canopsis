db.permission.updateOne({_id: "listalarm_pbehavior"}, {
    $set: {
        "api_permissions.api_pbehaviortype": 4,
        "api_permissions.api_pbehaviorreason": 4,
        "api_permissions.api_pbehaviorexception": 4,
    }
});

db.role.updateMany({"permissions.listalarm_pbehavior": 1}, [
    {
        $set: {
            "permissions.api_pbehaviortype": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_pbehaviortype", 0]}}]},
            "permissions.api_pbehaviorreason": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_pbehaviorreason", 0]}}]},
            "permissions.api_pbehaviorexception": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_pbehaviorexception", 0]}}]},

        }
    }
]);

db.permission.updateOne({_id: "listalarm_fastPbehavior"}, {
    $set: {
        "api_permissions.api_pbehavior": 9,
        "api_permissions.api_pbehaviortype": 4,
        "api_permissions.api_pbehaviorreason": 4,
    }
});

db.role.updateMany({"permissions.listalarm_fastPbehavior": 1}, [
    {
        $set: {
            "permissions.api_pbehavior": {$bitOr: [9, {$toInt: {$ifNull: ["$permissions.api_pbehavior", 0]}}]},
            "permissions.api_pbehaviortype": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_pbehaviortype", 0]}}]},
            "permissions.api_pbehaviorreason": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_pbehaviorreason", 0]}}]},

        }
    }
]);
