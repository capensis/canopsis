(function () {
    db.periodical_alarm.find({
        "v.state._t": "changestate",
        "v.state.val": 0,
        "v.status.val": 1,
    }).forEach(function (doc) {
        var statusStep = {
            "_t": "statusdec",
            "t": doc.v.state.t,
            "a": doc.v.state.a,
            "m": doc.v.state.m,
            "val": 0,
            "initiator": doc.v.state.initiator,
        };
        db.periodical_alarm.updateOne({_id: doc._id}, {
            $set: {
                "v.status": statusStep,
            },
            $push: {
                "v.steps": statusStep,
            }
        });
    });
})();
