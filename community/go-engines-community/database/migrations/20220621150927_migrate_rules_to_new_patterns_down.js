// down script can be applied only if no rule has been updated by API
db.default_entities.updateMany({type: "service"}, {
    $rename: {
        old_entity_patterns: "entity_patterns",
    },
    $unset: {
        entity_pattern: "",
    },
});

db.kpi_filter.updateMany({}, {
    $rename: {
        old_entity_patterns: "entity_patterns",
    },
    $unset: {
        entity_pattern: "",
    },
});

for (var collectionName of ["idle_rule", "dynamic_infos", "instruction", "resolve_rule", "flapping_rule"]) {
    var collection = db.getCollection(collectionName);
    collection.updateMany({}, {
        $rename: {
            old_entity_patterns: "entity_patterns",
            old_alarm_patterns: "alarm_patterns",
        },
        $unset: {
            entity_pattern: "",
            alarm_pattern: "",
        },
    });
}

db.action_scenario.find().forEach(function (doc) {
    var newActions = [];

    doc.actions.forEach(function (action) {
        var newAction = action;

        if (newAction.old_entity_patterns) {
            newAction.entity_patterns = newAction.old_entity_patterns;
        }
        if (newAction.old_alarm_patterns) {
            newAction.alarm_patterns = newAction.old_alarm_patterns;
        }

        delete newAction.entity_pattern;
        delete newAction.alarm_pattern;

        newActions.push(newAction);
    });

    db.action_scenario.updateOne({_id: doc._id}, {$set: {actions: newActions}});
});

db.eventfilter.updateMany({}, {
    $rename: {
        old_patterns: "patterns",
    },
    $unset: {
        event_pattern: "",
        entity_pattern: "",
    },
});

db.meta_alarm_rules.updateMany({}, {
    $rename: {
        "old_alarm_patterns": "config.alarm_patterns",
        "old_entity_patterns": "config.entity_patterns",
        "old_total_entity_patterns": "config.total_entity_patterns",
        "old_event_patterns": "config.event_patterns",
    },
    $unset: {
        alarm_pattern: "",
        entity_pattern: "",
        total_entity_pattern: "",
    },
});
