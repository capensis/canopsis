db.share_token.dropIndex("value_1");

db.default_rights.deleteMany({
    _id: {
        $in: [
            "api_share_token",
            "models_shareToken",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_share_token": "",
        "rights.models_shareToken": "",
    }
});
