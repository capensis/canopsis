db.default_rights.deleteMany({
    _id: {
        $in: ["models_remediationStatistic"],
    }
});
db.default_rights.updateOne({_id: "admin"}, {
    $unset: {
        "rights.models_remediationStatistic": "",
    }
});
