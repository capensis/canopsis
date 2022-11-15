db.instruction_execution.aggregate([
    {
        $lookup: {
            from: "instruction",
            localField: "instruction",
            foreignField: "_id",
            as: "instruction"
        }
    },
    {$unwind: "$instruction"},
]).forEach(function (doc) {
    var set = {
        name: doc.instruction.name,
        description: doc.instruction.description,
        type: doc.instruction.type,
        timeout_after_execution: doc.instruction.timeout_after_execution,
        created_at: doc.started_at,
        priority: doc.instruction.priority,
    };
    if (doc.instruction.type == 0 && doc.instruction_modified_on >= doc.instruction.last_modified && doc.step_history && doc.instruction.steps) {
        var updated = true
        doc.step_history.forEach(function (step, stepIndex) {
            var instructionStep = doc.instruction.steps[stepIndex];
            if (!instructionStep) {
                updated = false;
                return;
            }

            step.name = instructionStep.name;
            step.stop_on_fail = instructionStep.stop_on_fail;
            step.endpoint = instructionStep.endpoint;

            if (!step.operation_history || !instructionStep.operations) {
                updated = false;
                return;
            }

            step.operation_history.forEach(function (op, opIndex) {
                var instructionOp = instructionStep.operations[opIndex];
                if (!instructionOp) {
                    updated = false;
                    return;
                }

                op.name = instructionOp.name;
                op.time_to_complete = instructionOp.time_to_complete;
                op.jobs = instructionOp.jobs;
                op.files = instructionOp.files;
            });
        });

        if (updated) {
            set.step_history = doc.step_history;
        }
    }

    db.instruction_execution.updateOne({_id: doc._id}, {$set: set});
});
