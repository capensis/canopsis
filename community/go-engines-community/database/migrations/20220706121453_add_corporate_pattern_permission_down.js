db.default_rights.deleteMany({
    _id: {
        $in: ["api_corporate_pattern", "models_profile_corporatePattern"]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_corporate_pattern": "",
        "rights.models_profile_corporatePattern": "",
    }
});
