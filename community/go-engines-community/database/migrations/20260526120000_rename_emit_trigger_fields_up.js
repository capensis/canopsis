db.action_scenario.find({actions: {$exists: true, $ne: []}}).forEach(function (doc) {
    var rename = {};
    for (var index = 0; index < doc.actions.length; index++) {
        rename["actions." + index + ".emit_trigger"] = "actions." + index + ".emit_trigger_success";
    }
    db.action_scenario.updateOne({_id: doc._id}, {$rename: rename});
});

db.declare_ticket_rule.updateMany({}, {$rename: {emit_trigger: "emit_trigger_success"}});
db.webhook_history.updateMany({}, {$rename: {emit_trigger: "emit_trigger_success"}});
