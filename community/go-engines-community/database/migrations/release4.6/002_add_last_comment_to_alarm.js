(function () {
    var filter = {"v.steps._t": "comment"};
    var updatePipeline = [
        {
            $set: {
                "v.last_comment": {
                    "$arrayElemAt": [
                        {
                            "$filter": {
                                "input": "$v.steps",
                                "cond": {
                                    "$eq": ["$$this._t", "comment"],
                                },
                            }
                        },
                        -1,
                    ],
                }
            }
        },
    ];

    db.periodical_alarm.updateMany(filter, updatePipeline);
    db.resolved_alarms.updateMany(filter, updatePipeline);
})();
