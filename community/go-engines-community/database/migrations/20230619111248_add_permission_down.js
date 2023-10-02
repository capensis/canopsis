db.default_rights.deleteMany({
    _id: "listalarm_unCancel",
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.listalarm_unCancel": "",
    }
});
