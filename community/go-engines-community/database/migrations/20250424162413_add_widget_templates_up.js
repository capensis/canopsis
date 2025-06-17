// Available global functions:
// genID returns a new string UUID
// isInt checks if a value is integer
// toInt transforms value to integer

const now = Math.ceil((new Date()).getTime() / 1000);

if (db.widget_templates.countDocuments({type: "alarm_quick_actions"}) === 0) {
    db.widget_templates.insertOne({
        _id: genID(),
        title: "Default quick actions",
        type: "alarm_quick_actions",
        actions: ["ack", "fastAck", "cancel"],
        author: "root",
        created: now,
        updated: now
    });
}

if (db.widget_templates.countDocuments({type: "alarm_mass_quick_actions"}) === 0) {
    db.widget_templates.insertOne({
        _id: genID(),
        title: "Default quick massive actions",
        type: "alarm_mass_quick_actions",
        actions: ["ack", "fastAck", "cancel"],
        author: "root",
        created: now,
        updated: now
    });
}
