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
    db.instruction_execution.updateOne({_id: doc._id}, {
        $set: {
            name: doc.instruction.name,
            type: doc.instruction.type,
            timeout_after_execution: doc.instruction.timeout_after_execution,
            created_at: doc.started_at,
            priority: doc.instruction.priority,
        }
    });
});
