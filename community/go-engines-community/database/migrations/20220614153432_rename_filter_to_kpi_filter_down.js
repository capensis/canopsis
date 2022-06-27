db.kpi_filter.renameCollection("filter");

db.default_rights.insertOne({
    _id: "api_filter",
    crecord_name: "api_filter",
    crecord_type: "action",
    desc: "Filters",
    type: "CRUD"
});
db.default_rights.updateMany(
    {"rights.api_kpi_filter": {$ne: null}},
    [
        {$set: {"rights.api_filter": "$rights.api_kpi_filter"}},
        {$unset: "rights.api_kpi_filter"}
    ]
);
db.default_rights.deleteOne({_id: "api_kpi_filter"});
