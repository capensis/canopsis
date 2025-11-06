db.permission.updateMany({_id: {$in: ["models_exploitation_linkRule", "models_exploitation_eventFilter"]}}, {
    $unset: {
        "api_permissions.api_external_data_table": ""
    }
});
