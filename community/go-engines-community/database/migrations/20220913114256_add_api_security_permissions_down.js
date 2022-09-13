db.default_rights.deleteMany({
    _id: {
        $in: ["api_security_read", "api_security_update"]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_security_read": "",
        "rights.api_security_update": "",
    }
});
