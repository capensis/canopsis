(function () {
    function genID() {
        return UUID().toString().split('"')[1];
    }

    var now = Math.ceil((new Date()).getTime() / 1000);
    // Get priority
    var priority = -1;
    db.action_scenario.find().sort({priority: -1}).limit(1).forEach(function (doc) {
        priority = doc.priority
    })

    db.default_action.find().sort({priority: 1}).forEach(function (doc) {
        // Get priority
        priority += 1;
        // Generate name
        var name = "action-" + doc._id
        if (doc.hook.event_patterns && doc.hook.event_patterns.length > 0) {
            name += " (event_patterns cannot be migrated " + JSON.stringify(doc.hook.event_patterns) + ")";
        }
        if (doc.fields && doc.fields.length > 0) {
            name += " (fields cannot be migrated " + JSON.stringify(doc.fields) + ")";
        }
        if (doc.regex && doc.regex !== "") {
            name += " (regex cannot be migrated " + JSON.stringify(doc.regex) + ")";
        }
        // Get parameters
        var parameters = doc.parameters;
        if (parameters.message) {
            parameters.output = parameters.message;
            delete parameters.message;
        }
        switch (doc.type) {
            case "snooze":
                parameters.duration = {
                    seconds: parameters.duration,
                    unit: "s"
                };
                break
            case "declareticket":
                return;
        }
        // Get disable_during_periods
        var disable_during_periods = [];
        if (doc.disable_during_periods && doc.disable_during_periods.length > 0) {
            disable_during_periods = doc.disable_during_periods;
        }
        // Get triggers
        var triggers = doc.hook.triggers;
        // Get delay
        var delay = null
        if (doc.delay && doc.delay !== "") {
            var val = parseInt(doc.delay);
            var unit = doc.delay.replace(val.toString(), "");
            var seconds = val;
            switch (unit) {
                case "s":
                    seconds = val;
                    break
                case "m":
                    seconds = val * 60;
                    break
                case "h":
                    seconds = val * 60 * 60;
                    break
            }

            if (seconds > 0) {
                delay = {
                    seconds: seconds,
                    unit: unit,
                }
            }
        }
        // Get enabled
        var enabled = true;
        if (doc.enabled === false || doc.enabled === true) {
            enabled = doc.enabled;
        }
        // Get created
        var created = now;
        if (doc.created) {
            created = doc.creation_date;
        }
        // Get action
        var action = {
            type: doc.type,
            parameters: parameters,
            drop_scenario_if_not_matched: false,
            emit_trigger: false
        };
        if (doc.hook.alarm_patterns) {
            action.alarm_patterns = doc.hook.alarm_patterns;
        }
        if (doc.hook.entity_patterns) {
            action.entity_patterns = doc.hook.entity_patterns;
        }
        // Get author
        var author = "root"
        if (doc.parameters.author) {
            author = doc.parameters.author;
        }
        // Insert scenario
        db.action_scenario.insertOne({
            _id: genID(),
            name: name,
            author: author,
            enabled: enabled,
            disable_during_periods: disable_during_periods,
            triggers: triggers,
            delay: delay,
            actions: [action],
            priority: priority,
            created: created,
            updated: now,
        });
    });
})();
