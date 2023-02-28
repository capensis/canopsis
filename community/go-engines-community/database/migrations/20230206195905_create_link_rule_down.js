db.link_rule.drop();

db.default_rights.deleteMany({
    _id: {
        $in: [
            "api_link_rule",
            "models_exploitation_linkRule",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_link_rule": "",
        "rights.models_exploitation_linkRule": "",
    }
});
