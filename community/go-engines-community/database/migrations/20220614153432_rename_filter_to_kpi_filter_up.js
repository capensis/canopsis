if (!db.default_rights.findOne({_id: "api_kpi_filter"})) {
    db.filter.renameCollection("kpi_filter");

    db.default_rights.insertOne({
        _id: "api_kpi_filter",
        crecord_name: "api_kpi_filter",
        crecord_type: "action",
        desc: "KPI filters",
        type: "CRUD"
    });
    db.default_rights.updateMany(
        {"rights.api_filter": {$ne: null}},
        [
            {$set: {"rights.api_kpi_filter": "$rights.api_filter"}},
            {$unset: "rights.api_filter"}
        ]
    );
    db.default_rights.deleteOne({_id: "api_filter"});
}

db.createCollection("kpi_filter");
