(function () {
    var pipeline = [
        {
            $lookup: {
                from: "periodical_alarm",
                localField: "alarm",
                foreignField: "_id",
                as: "alarm",
            }
        },
        {$unwind: {path: "$alarm", preserveNullAndEmptyArrays: true}},
        {
            $lookup: {
                from: "instruction",
                localField: "instruction",
                foreignField: "_id",
                as: "instruction",
            }
        },
        {$unwind: {path: "$instruction", preserveNullAndEmptyArrays: true}},
    ];
    var cursor = db.instruction_execution.aggregate(pipeline);

    while (cursor.hasNext()) {
        var doc = cursor.next();
        var set = {};

        if (!doc.started_at) {
            if (doc.step_history && doc.step_history.length > 0) {
                var firstStep = doc.step_history[0];
                if (firstStep.operation_history && firstStep.operation_history.length > 0) {
                    var firstOperation = firstStep.operation_history[0];
                    if (firstOperation.started_at) {
                        doc.started_at = firstOperation.started_at;
                        set["started_at"] = doc.started_at;
                    }
                }
            }
        }

        if (!doc.completed_at) {
            if (doc.step_history && doc.step_history.length > 0) {
                var lastStep = doc.step_history[doc.step_history.length - 1];
                if (lastStep.operation_history && lastStep.operation_history.length > 0) {
                    var lastOperation = lastStep.operation_history[lastStep.operation_history.length - 1];
                    if (lastOperation.completed_at) {
                        doc.completed_at = lastOperation.completed_at;
                        set["completed_at"] = doc.completed_at;
                    }
                }
            }
        }

        if (!doc.complete_time && doc.status === 2) {
            var completeTime = 0;
            if (doc.step_history) {
                doc.step_history.forEach(function (step) {
                    if (step.operation_history) {
                        step.operation_history.forEach(function (operation) {
                            if (operation.completed_at && operation.started_at) {
                                completeTime += operation.completed_at - operation.started_at;
                            }
                        });
                    }
                });

                if (completeTime > 0) {
                    set["complete_time"] = completeTime;
                }
            }
        }

        if (!doc.alarm_state && doc.alarm && doc.alarm.v.steps) {
            var startIndex = -1;
            doc.alarm.v.steps.forEach(function (step, i) {
                if (startIndex >= 0) {
                    return;
                }
                if (step._t === "instructionstart" && step.t >= doc.started_at && step.t <= doc.completed_at) {
                    startIndex = i;
                }
            });
            var state = -1;
            if (startIndex >= 0) {
                doc.alarm.v.steps.forEach(function (step, i) {
                    if (i >= startIndex) {
                        return;
                    }
                    if (step._t === "statedec" || step._t === "stateinc") {
                        state = step.val;
                    }
                });
            }

            if (state >= 0) {
                set["alarm_state"] = state;

                if (doc.status === 2) {
                    var completeIndex = -1;
                    var t = -1;
                    doc.alarm.v.steps.forEach(function (step, i) {
                        if (completeIndex >= 0 || i <= startIndex) {
                            return;
                        }
                        if (step._t === "instructioncomplete") {
                            completeIndex = i;
                            t = step.t + 60;
                        }
                    });
                    var resState = -1;
                    if (completeIndex >= 0) {
                        doc.alarm.v.steps.forEach(function (step, i) {
                            if (i <= completeIndex || resState >= 0) {
                                return;
                            }
                            if (step.t <= t && (step._t === "statedec" || step._t === "stateinc")) {
                                resState = step.val;
                            }
                        });
                    }
                    if (resState < 0) {
                        resState = state;
                    }

                    set["result_alarm_state"] = resState;
                }
            }
        }

        if (!doc.instruction_modified_on && doc.instruction) {
            if (doc.started_at < doc.instruction.last_modified && doc.instruction.created) {
                set["instruction_modified_on"] = doc.instruction.created;
            } else {
                set["instruction_modified_on"] = doc.instruction.last_modified;
            }
        }

        if (Object.keys(set).length > 0) {
            db.instruction_execution.updateOne({_id: doc._id}, {"$set": set});
        }
    }
})();
