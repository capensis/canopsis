db.default_rights.deleteMany({
    _id: {
        $in: ["api_techmetrics"]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_techmetrics": "",
    }
});
