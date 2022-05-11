(function () {
    function genID() {
        return UUID().toString().split('"')[1];
    }
    var now = Math.ceil((new Date()).getTime() / 1000);

    db.instruction.find().forEach(function (doc) {
        if (!doc.created) {doc.created = now}
        var modStats = [
            {
                _id: genID(),
                instruction: doc._id,
                date: doc.created,
            }
        ];
        if (doc.last_modified !== doc.created) {
            modStats.push({
                _id: genID(),
                instruction: doc._id,
                date: doc.last_modified,
            })
        }

        db.instruction_mod_stats.insertMany(modStats);
    });

    var cursor = db.instruction_execution.aggregate([
        {$match: {status: 2}},
        {
            $group: {
                _id: {
                    "instruction_modified_on": "$instruction_modified_on",
                    "instruction": "$instruction",
                },
                date: {$first: "$instruction_modified_on"},
                instruction: {$first: "$instruction"},
                execution_count: {$sum: 1},
                complete_time: {$sum: "$complete_time"},
                init_critical: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$alarm_state", 3]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                init_major: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$alarm_state", 2]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                init_minor: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$alarm_state", 1]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                res_critical: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$result_alarm_state", 3]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                res_major: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$result_alarm_state", 2]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                res_minor: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$result_alarm_state", 1]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
                res_ok: {
                    $sum: {
                        $cond: {
                            if: {$eq: ["$result_alarm_state", 0]},
                            then: 1,
                            else: 0,
                        }
                    }
                },
            }
        }
    ]);
    while (cursor.hasNext()) {
        var doc = cursor.next();

        if (doc.execution_count === 0) {
            continue;
        }

        db.instruction_mod_stats.updateOne({instruction: doc.instruction, date: doc.date}, {
            $set: {avg_complete_time: Math.round(doc.complete_time / doc.execution_count)},
            $inc: {
                execution_count: doc.execution_count,
                init_critical: doc.init_critical,
                init_major: doc.init_major,
                init_minor: doc.init_minor,
                res_critical: doc.res_critical,
                res_major: doc.res_major,
                res_minor: doc.res_minor,
                res_ok: doc.res_ok,
            }
        });
    }
})();
