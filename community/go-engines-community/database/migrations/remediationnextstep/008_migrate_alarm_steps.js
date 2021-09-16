(function () {
    var stepTypes = [
        "instructionstart",
        "instructionpause",
        "instructionresume",
        "instructioncomplete",
        "instructionfail",
        "instructionabort",
        "instructionjobstart",
        "instructionjobcomplete",
        "instructionjobabort",
        "instructionjobfail",
    ];
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
        if (!doc.instruction) {
            return;
        }

        var output = "Instruction " + doc.instruction.name;

        if (!doc.alarm || !doc.alarm.v || !doc.alarm.v.steps) {
            return;
        }

        var steps = [];
        doc.alarm.v.steps.forEach(function (step) {
            if (step._t && step.m && stepTypes.includes(step._t) && step.m.startsWith(output)) {
                step.exec = doc._id;
            }

            steps.push(step);
        });

        db.periodical_alarm.updateOne({_id: doc.alarm._id}, {
            "$set": {
                "v.steps": steps,
            }
        });
    }
})();
