db.action_scenario.find({actions: {$exists: true, $ne: []}}).forEach(function (doc) {
    var set = {};
    var unset = {};

    for (var index = 0; index < doc.actions.length; index++) {
        var action = doc.actions[index];

        if (action.emit_trigger_success !== undefined) {
            set["actions." + index + ".emit_trigger"] = action.emit_trigger_success;
            unset["actions." + index + ".emit_trigger_success"] = "";
        }

        unset["actions." + index + ".emit_trigger_fail"] = "";
    }

    if (Object.keys(set).length > 0) {
        db.action_scenario.updateOne({_id: doc._id}, {$set: set});
    }

    if (Object.keys(unset).length > 0) {
        db.action_scenario.updateOne({_id: doc._id}, {$unset: unset});
    }
});

db.declare_ticket_rule.updateMany({}, {$rename: {emit_trigger_success: "emit_trigger"}});
db.declare_ticket_rule.updateMany({}, {$unset: {emit_trigger_fail: ""}});

db.webhook_history.updateMany({}, {$rename: {emit_trigger_success: "emit_trigger"}});
db.webhook_history.updateMany({}, {$unset: {emit_trigger_fail: ""}});
