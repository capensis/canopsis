db.default_rights.deleteOne({_id: "models_remediationStatistic"});
db.default_rights.updateOne({_id: "admin"}, {
    $unset: {
        "rights.models_remediationStatistic": "",
    }
});
