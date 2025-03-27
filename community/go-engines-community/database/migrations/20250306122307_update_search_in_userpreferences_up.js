// Available global functions:
// genID returns a new string UUID
// isInt checks if a value is integer
// toInt transforms value to integer

db.userpreferences.find({"content.searches": {$ne: null}}).forEach(function (doc) {
    let searches = [];
    let updated = false;
    for (const search of doc.content.searches) {
        if (!search._id) {
            updated = true;
            searches.push({
                _id: genID(),
                ...search,
            });
        } else {
            searches.push(search);
        }
    }

    if (updated) {
        db.userpreferences.updateOne({_id: doc._id}, {$set: {"content.searches": searches}});
    }
});
