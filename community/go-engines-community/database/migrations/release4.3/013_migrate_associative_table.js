(function () {
    db.default_associativetable.find().forEach(function (doc) {
        db.default_associativetable.updateOne({name: doc.name}, {
            $set: {
                content: {
                    val: doc.content,
                }
            }
        })
    });
})();