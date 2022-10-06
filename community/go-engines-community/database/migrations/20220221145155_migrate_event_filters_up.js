db.eventfilter.find({
    "type": "enrichment"
}).forEach(function (doc) {
    if (doc.actions === undefined) {
        return
    }

    doc.actions.forEach(function (action) {
        if (action.from !== undefined) {
            action.value = action.from
            delete action.from
        }

        if (action.to !== undefined) {
            action.name = action.to
            delete action.to
        }
    })

    db.eventfilter.updateOne({_id: doc._id}, {
        $set: {
            config: {
                actions: doc.actions,
                on_success: doc.on_success,
                on_failure: doc.on_failure
            }
        },
        $unset: {
            actions: "",
            on_success: "",
            on_failure: ""
        }
    });
});

db.eventfilter.createIndex({
    priority: 1,
}, {name: "priority_1"});
