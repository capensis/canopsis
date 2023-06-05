var set = {};
db.views.find().forEach(function (doc) {
    set["rights." + doc._id] = {
        checksum: 7
    };
});
db.view_playlist.find().forEach(function (doc) {
    set["rights." + doc._id] = {
        checksum: 7
    };
});
if (Object.keys(set).length > 0) {
    db.default_rights.updateOne({_id: "admin"}, {$set: set});
}
