db.widgets.updateMany(
    {
        "type": "Context",
        "parameters.serviceDependenciesColumns": null,
    },
    {
        $set: {
            "parameters.serviceDependenciesColumns": [
                {
                    "value": "name",
                    "label": "Name"
                },
                {
                    "value": "type",
                    "label": "Type"
                }
            ],
        },
    },
);
