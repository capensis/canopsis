(function () {
    db.default_associativetable.find({ "name": { $ne: "link_builders_settings"}}).forEach(function (doc) {
        db.default_associativetable.updateOne({name: doc.name}, {
            $set: {
                content: {
                    val: doc.content,
                }
            }
        })
    });
})();
