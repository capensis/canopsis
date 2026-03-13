db.eventfilter_failure.aggregate([
    {
        $match: {rule_updated: null}
    },
    {
        $group: {_id: "$rule"}
    },
    {
        $lookup: {
            from: "eventfilter",
            foreignField: "_id",
            localField: "_id",
            as: "rule",
        }
    },
    {
        $unwind: "$rule"
    },
    {
        $project: {
            updated: "$rule.updated"
        }
    }
]).forEach(function (doc) {
    db.eventfilter_failure.updateMany(
        {rule: doc._id, rule_updated: null},
        {$set: {rule_updated: doc.updated}}
    );
});

db.eventfilter_failure.dropIndex("rule_1");
db.eventfilter_failure.createIndex({rule: 1, rule_updated: 1}, {name: "rule_1_rule_updated_1"});
db.eventfilter_failure.createIndex({t: 1}, {name: "t_1"});

db.user_notification.aggregate([
    {
        $match: {
            type: 1,
            "rule.updated": null
        }
    },
    {
        $group: {_id: "$rule._id"}
    },
    {
        $lookup: {
            from: "eventfilter",
            foreignField: "_id",
            localField: "_id",
            as: "rule",
        }
    },
    {
        $unwind: "$rule"
    },
    {
        $project: {
            updated: "$rule.updated"
        }
    }
]).forEach(function (doc) {
    db.user_notification.updateMany(
        {type: 1, "rule._id": doc._id, "rule.updated": null},
        {$set: {"rule.updated": doc.updated}}
    );
});
