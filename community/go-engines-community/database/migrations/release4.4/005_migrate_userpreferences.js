(function () {
    var cursor = db.userpreferences.aggregate([
        {
            $match: {
                widgetXtype: "AlarmsList",
                "widget_preferences.remediationInstructionsFilters": {$ne: null}
            }
        },
        {
            $lookup: {
                from: "instruction",
                localField: "widget_preferences.remediationInstructionsFilters.instructions",
                foreignField: "name",
                as: "instructions",
            }
        }
    ]);

    while (cursor.hasNext()) {
        var doc = cursor.next();
        var instructionsByName = {};
        doc.instructions.forEach(function (instruction) {
            instructionsByName[instruction.name] = instruction;
        });

        var filters = [];
        doc.widget_preferences.remediationInstructionsFilters.forEach(function (filter) {
            var newInstructions = [];
            filter.instructions.forEach(function (name) {
                var instruction = instructionsByName[name];
                if (instruction !== undefined) {
                    newInstructions.push({
                        _id: instruction._id,
                        name: instruction.name,
                        type: instruction.type,
                    })
                }
            });

            filter.instructions = newInstructions;
            filters.push(filter);
        });

        db.userpreferences.updateOne({_id: doc._id}, {
            $set: {
                "widget_preferences.remediationInstructionsFilters": filters,
            }
        });
    }
})();