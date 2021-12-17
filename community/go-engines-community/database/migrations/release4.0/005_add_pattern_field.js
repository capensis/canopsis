(function () {
    function extractDeepestValue(docPatterns, valuePatterns){
        if (docPatterns == null){
            return;
        }
        switch (typeof docPatterns){
            case "object":
                for (var key of Object.keys(docPatterns)) {
                    extractDeepestValue(docPatterns[key], valuePatterns);
                }
                break;
            case "boolean":
            case "number":
            case "bigint":
            case "string":
                valuePatterns.push(docPatterns);
        }
    }

    db.dynamic_infos.find().forEach(function (doc){
        const patterns = [];
        extractDeepestValue(doc.entity_patterns, patterns);
        extractDeepestValue(doc.alarm_patterns, patterns);
        db.dynamic_infos.update(
            { _id: doc._id },
            { "$set": { "pattern": patterns } }
        );
    });
})();