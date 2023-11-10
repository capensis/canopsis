// items in kpi_filter, meta_alarm_rules, widget_filter, resolve_rule, flapping_rule cannot be disabled

db.default_entities.aggregate([
    {
        $match: {
            type: "service",
            entity_pattern: {$in: [null, []]},
            old_entity_patterns: {$ne: null},
            enabled: true,
        }
    },
    {
        $lookup: {
            from: "periodical_alarm",
            let: {entity: "$_id"},
            pipeline: [
                {
                    $match: {
                        $and: [
                            {"v.resolved": null},
                            {$expr: {$eq: ["$d", "$$entity"]}}
                        ]
                    }
                }
            ],
            as: "alarm",
        }
    },
    {$unwind: {path: "$alarm", preserveNullAndEmptyArrays: true}},
    {
        $project: {
            alarm: 1
        }
    }
]).forEach(function (doc) {
    db.default_entities.updateOne({_id: doc._id}, {$set: {enabled: false}});
    if (doc.alarm) {
        var now = Math.ceil((new Date()).getTime() / 1000);
        db.periodical_alarm.updateOne({_id: doc.alarm._id}, {$set: {"v.resolved": now}});
        doc.alarm.v.resolved = now;
        db.resolved_alarms.updateOne({_id: doc.alarm._id}, {$set: doc.alarm}, {upsert: true});
    }
});

for (var collectionName of ["idle_rule", "dynamic_infos", "instruction"]) {
    var collection = db.getCollection(collectionName);
    collection.updateMany(
        {
            $or: [
                {
                    entity_pattern: {$in: [null, []]},
                    old_entity_patterns: {$ne: null},
                },
                {
                    alarm_pattern: {$in: [null, []]},
                    old_alarm_patterns: {$ne: null},
                },
            ],
            enabled: true,
        },
        {$set: {enabled: false}}
    );
}

db.action_scenario.find({enabled: true}).forEach(function (doc) {
    var update = false;
    doc.actions.forEach(function (action) {
        if (action.old_entity_patterns && (!action.entity_pattern || action.entity_pattern.length === 0) ||
            action.old_alarm_patterns && (!action.alarm_pattern || action.alarm_pattern.length === 0)) {
            update = true;
        }
    });

    if (update) {
        db.action_scenario.updateOne({_id: doc._id}, {$set: {enabled: false}});
    }
});

db.eventfilter.updateMany(
    {
        event_pattern: {$in: [null, []]},
        entity_pattern: {$in: [null, []]},
        old_patterns: {$ne: null},
        enabled: true,
    },
    {$set: {enabled: false}}
);

db.pbehavior.updateMany(
    {
        entity_pattern: {$in: [null, []]},
        old_mongo_query: {$ne: null},
        enabled: true,
    },
    {$set: {enabled: false}}
);
