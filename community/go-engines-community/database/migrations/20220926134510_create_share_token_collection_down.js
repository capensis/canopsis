db.share_token.dropIndex("value_1");

db.default_rights.deleteOne({_id: "api_share_token"});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_share_token": "",
    }
});
