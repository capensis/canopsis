db.eventfilter.find({
    "type": "enrichment"
}).forEach(function (doc) {
    doc.actions.forEach(function (action) {
        if (action.value !== undefined) {
            action.from = action.value
            delete action.value
        }

        if (action.name !== undefined) {
            action.to = action.name
            delete action.name
        }
    })

    db.eventfilter.updateOne({_id: doc._id}, {
        $set: {
            actions: doc.config.actions,
            on_success: doc.config.on_success,
            on_failure: doc.config.on_failure,
        },
        $unset: {
            "config.actions": "",
            "config.on_success": "",
            "config.on_failure": ""
        }
    });
});

db.eventfilter.dropIndex("priority_1");
