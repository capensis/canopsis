db.action_scenario.find({actions: {$exists: true, $ne: []}}).forEach(function (doc) {
    var rename = {};
    var unset = {};
    for (var index = 0; index < doc.actions.length; index++) {
        rename["actions." + index + ".emit_trigger_success"] = "actions." + index + ".emit_trigger";
        unset["actions." + index + ".emit_trigger_fail"] = "";
    }
    db.action_scenario.updateOne({_id: doc._id}, {$rename: rename});
    db.action_scenario.updateOne({_id: doc._id}, {$unset: unset});
});

db.declare_ticket_rule.updateMany({}, {$rename: {emit_trigger_success: "emit_trigger"}});
db.declare_ticket_rule.updateMany({}, {$unset: {emit_trigger_fail: ""}});

db.webhook_history.updateMany({}, {$rename: {emit_trigger_success: "emit_trigger"}});
db.webhook_history.updateMany({}, {$unset: {emit_trigger_fail: ""}});
