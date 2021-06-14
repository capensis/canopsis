(function () {
    db.default_rights.find({
        crecord_type: "action",
        $or: [
            {crecord_name: null},
            {crecord_name: ""}
        ],
    }).forEach(function (doc) {
        var view = db.views.findOne({_id: doc._id});
        if (view) {
            db.default_rights.updateOne({_id: doc._id}, {
                $set: {
                    crecord_name: doc._id,
                    desc: "Rights on view : " + view.title
                }
            });
        } else {
            var playlist = db.view_playlist.findOne({_id: doc._id});
            if (playlist) {
                db.default_rights.updateOne({_id: doc._id}, {
                    $set: {
                        crecord_name: doc._id,
                        desc: "Rights on playlist : " + playlist.name,
                        type: "RW"
                    }
                });
            } else {
                db.default_rights.updateOne({_id: doc._id}, {
                    $set: {
                        crecord_name: doc._id
                    }
                });
            }
        }
    });
})();
