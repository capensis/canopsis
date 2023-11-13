db.state_settings.dropIndex("title_1");

db.state_settings.deleteOne({_id: "service"});
db.state_settings.updateOne({_id: "junit"}, {
    $set: {
        type: "junit",
    },
    $unset: {
        title: "",
        on_top: "",
        enabled: ""
    }
});

db.permission.deleteMany({_id: "models_stateSetting"});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.models_stateSetting": "",
        }
    },
);

db.permission.updateOne({_id:"api_state_settings"}, {$unset: {type: ""}});
db.role.updateMany({"permissions.api_state_settings": {$ne: null}}, {$set: {"permissions.api_state_settings": 1}});

