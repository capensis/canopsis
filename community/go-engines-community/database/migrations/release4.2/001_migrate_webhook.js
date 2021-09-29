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

    db.webhooks.find().forEach(function (doc) {
        // Get priority
        priority += 1;
        // Get declare_ticket
        var declare_ticket = null;
        if (doc.declare_ticket && Object.keys(doc.declare_ticket).length > 0) {
            declare_ticket = doc.declare_ticket;
            if (declare_ticket.regexp) {
                declare_ticket.is_regexp = declare_ticket.regexp;
                delete declare_ticket.regexp;
            }
        }
        // Generate name
        var name = "webhook-" + doc._id;
        if (doc.hook.event_patterns && doc.hook.event_patterns.length > 0) {
            name += " (event_patterns cannot be migrated " + JSON.stringify(doc.hook.event_patterns) + ")";
        }
        // Get retry
        var retry_count = null;
        var retry_delay = null;
        if (doc.retry && Object.keys(doc.retry).length > 0) {
            retry_count = doc.retry.count;
            var val = doc.retry.delay;
            var unit = doc.retry.unit;
            if (val > 0) {
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

                retry_delay = {
                    seconds: seconds,
                    unit: unit,
                };
            }
        }
        // Get disable_during_periods
        var disable_during_periods = [];
        if (doc.disable_during_periods && doc.disable_during_periods.length > 0) {
            disable_during_periods = doc.disable_during_periods;
        }
        // Get triggers
        var triggers = doc.hook.triggers;
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
            type: 'webhook',
            parameters: {
                declare_ticket: declare_ticket,
                retry_count: retry_count,
                retry_delay: retry_delay,
                request: doc.request,
            },
            drop_scenario_if_not_matched: false,
            emit_trigger: true
        };
        if (doc.hook.alarm_patterns) {
            action.alarm_patterns = doc.hook.alarm_patterns;
        }
        if (doc.hook.entity_patterns) {
            action.entity_patterns = doc.hook.entity_patterns;
        }
        // Get author
        var author = "root"
        if (doc.author) {
            author = doc.author;
        }
        // Insert scenario
        db.action_scenario.insertOne({
            _id: genID(),
            name: name,
            author: author,
            enabled: enabled,
            disable_during_periods: disable_during_periods,
            triggers: triggers,
            actions: [action],
            priority: priority,
            created: created,
            updated: now,
        });
    });
})();