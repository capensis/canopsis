db.action_scenario.find({actions: {$exists: true, $ne: []}}).forEach(function (doc) {
    var set = {};
    var unset = {};

    for (var index = 0; index < doc.actions.length; index++) {
        var action = doc.actions[index];

        if (action.emit_trigger !== undefined) {
            set["actions." + index + ".emit_trigger_success"] = action.emit_trigger;
            unset["actions." + index + ".emit_trigger"] = "";
        }
    }

    if (Object.keys(set).length > 0) {
        db.action_scenario.updateOne(
            {
                _id: doc._id,
            },
            {
                $set: set,
                $unset: unset
            }
        );
    }
});

db.declare_ticket_rule.updateMany({}, {$rename: {emit_trigger: "emit_trigger_success"}});
db.webhook_history.updateMany({}, {$rename: {emit_trigger: "emit_trigger_success"}});
