(function () {
    db.eventfilter.find({patterns:{$elemMatch:{type:"watcher"}}}).forEach(function (doc) {
        doc.patterns.forEach(function (pattern) {
            if (pattern.type === "watcher") {
                pattern.type = "service"
            }
        })

        db.eventfilter.updateOne({_id: doc._id}, {
            "$set": {"patterns": doc.patterns},
        });
    })

    db.instruction.find({entity_patterns:{$elemMatch:{type:"watcher"}}}).forEach(function (doc) {
        doc.entity_patterns.forEach(function (pattern) {
            if (pattern.type === "watcher") {
                pattern.type = "service"
            }
        })

        db.instruction.updateOne({_id: doc._id}, {
            "$set": {"entity_patterns": doc.entity_patterns},
        });
    })

    db.action_scenario.find({actions:{$elemMatch:{entity_patterns:{$elemMatch:{type:"watcher"}}}}}).forEach(function (doc) {
        doc.actions.forEach(function (action) {
            action.entity_patterns.forEach(function (pattern) {
                if (pattern.type === "watcher") {
                    pattern.type = "service"
                }
            })
        })

        db.action_scenario.updateOne({_id: doc._id}, {
            "$set": {"actions": doc.actions},
        });
    })

    db.meta_alarm_rules.find({$or:[{"patterns.entity_patterns":{$elemMatch:{type:"watcher"}}}, {"config.entity_patterns":{$elemMatch:{type:"watcher"}}}]}).forEach(function (doc) {
        doc.patterns.entity_patterns.forEach(function (pattern) {
            if (pattern.type === "watcher") {
                pattern.type = "service"
            }
        })

        doc.config.entity_patterns.forEach(function (pattern) {
            if (pattern.type === "watcher") {
                pattern.type = "service"
            }
        })

        db.meta_alarm_rules.updateOne({_id: doc._id}, {
            "$set": {
                "patterns": doc.patterns,
                "config": doc.config
            },
        });
    })

    db.dynamic_infos.find({entity_patterns:{$elemMatch:{type:"watcher"}}}).forEach(function (doc) {
        doc.entity_patterns.forEach(function (pattern) {
            if (pattern.type === "watcher") {
                pattern.type = "service"
            }
        })

        db.dynamic_infos.updateOne({_id: doc._id}, {
            "$set": {"entity_patterns": doc.entity_patterns},
        });
    })

    db.pbehavior.find().forEach(function (doc) {
        db.pbehavior.updateOne({_id: doc._id}, {
            "$set": {"filter": doc.filter.replace(new RegExp("{\"type\":\"watcher\"}", 'g'), "{\"type\":\"service\"}")},
        });
    })
})();