db.default_rights.find({crecord_type: "action"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;

    db.permission.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.find({crecord_type: "role"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;
    doc.permissions = {};
    for (var permission of Object.keys(doc.rights)) {
        doc.permissions[permission] = doc.rights[permission].checksum;
    }
    delete doc.rights;

    db.role.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.find({crecord_type: "user"}).forEach(function (doc) {
    doc.name = doc.crecord_name;
    delete doc.crecord_name;
    delete doc.crecord_type;
    doc.password = doc.shadowpasswd;
    delete doc.shadowpasswd;
    if (doc.role) {
        doc.roles = [doc.role];
    } else {
        doc.roles = [];
    }
    delete doc.role;
    db.user.updateOne(
        {name: doc.name},
        {$setOnInsert: doc},
        {upsert: true},
    );
});

db.default_rights.drop();

if (!db.role_template.findOne({})) {
    db.role_template.insertMany([
        {
            "_id": genID(),
            "name": "Pilotes",
            "description": "Profil de pilotage",
            "permissions": {
                "serviceweather_entityManagePbehaviors": 1,
                "userStatistics_userFilter": 1,
                "alarmStatistics_userFilter": 1,
                "crudcontext_duplicate": 1,
                "models_remediationStatistic": 1,
                "numbers_listFilters": 1,
                "listalarm_editFilter": 1,
                "crudcontext_listFilters": 1,
                "listalarm_snoozeAlarm": 1,
                "serviceweather_entityLinks": 1,
                "crudcontext_edit": 1,
                "testingweather_alarmsList": 1,
                "numbers_interval": 1,
                "serviceweather_listFilters": 1,
                "userStatistics_listFilters": 1,
                "barchart_editFilter": 1,
                "listalarm_executeInstruction": 1,
                "api_alarm_read": 1,
                "listalarm_addRemediationInstructionsFilter": 1,
                "serviceweather_pbehaviorList": 1,
                "crudcontext_deletePbehavior": 1,
                "linechart_interval": 1,
                "api_file": 4,
                "ac79f382-3a21-4f77-868d-a42e9309fa60": 1,
                "listalarm_links": 1,
                "map_category": 1,
                "barchart_listFilters": 1,
                "piechart_listFilters": 1,
                "piechart_addFilter": 1,
                "api_execution": 1,
                "listalarm_history": 1,
                "listalarm_comment": 1,
                "serviceweather_entityPause": 1,
                "serviceweather_entityPlay": 1,
                "crudcontext_pbehavior": 1,
                "userStatistics_editFilter": 1,
                "alarmStatistics_listFilters": 1,
                "barchart_sampling": 1,
                "piechart_userFilter": 1,
                "api_pbehaviorexception": 4,
                "listalarm_assignTicketNumber": 1,
                "listalarm_manualMetaAlarmGroup": 1,
                "linechart_listFilters": 1,
                "linechart_addFilter": 1,
                "api_event": 1,
                "listalarm_declareanIncident": 1,
                "serviceweather_entityAssocTicket": 1,
                "map_userFilter": 1,
                "api_pbehaviortype": 4,
                "serviceweather_alarmsList": 1,
                "crudcontext_massEnable": 1,
                "common_variablesHelp": 1,
                "barchart_interval": 1,
                "linechart_editFilter": 1,
                "numbers_sampling": 1,
                "numbers_editFilter": 1,
                "api_associative_table": 4,
                "serviceweather_editFilter": 1,
                "alarmStatistics_interval": 1,
                "linechart_sampling": 1,
                "api_viewgroup": 4,
                "listalarm_category": 1,
                "crudcontext_listPbehavior": 1,
                "map_listFilters": 1,
                "userStatistics_interval": 1,
                "linechart_userFilter": 1,
                "listalarm_ack": 1,
                "listalarm_metaAlarmGroup": 1,
                "listalarm_userRemediationInstructionsFilter": 1,
                "map_editFilter": 1,
                "barchart_userFilter": 1,
                "api_playlist": 4,
                "api_entitycategory": 4,
                "api_entityservice": 4,
                "listalarm_userFilter": 1,
                "serviceweather_entityCancel": 1,
                "alarmStatistics_editFilter": 1,
                "numbers_addFilter": 1,
                "api_pbehaviorreason": 4,
                "crudcontext_massDisable": 1,
                "map_addFilter": 1,
                "serviceweather_entityAck": 1,
                "serviceweather_moreInfos": 1,
                "crudcontext_editFilter": 1,
                "crudcontext_addFilter": 1,
                "crudcontext_userFilter": 1,
                "piechart_editFilter": 1,
                "listalarm_cancelAck": 1,
                "listalarm_addFilter": 1,
                "listalarm_correlation": 1,
                "serviceweather_category": 1,
                "serviceweather_executeInstruction": 1,
                "counter_alarmsList": 1,
                "listalarm_changeState": 1,
                "api_junit": 4,
                "crudcontext_category": 1,
                "crudcontext_exportAsCsv": 1,
                "listalarm_removeAlarm": 1,
                "listalarm_listRemediationInstructionsFilters": 1,
                "listalarm_exportAsCsv": 1,
                "barchart_addFilter": 1,
                "api_alarm_update": 1,
                "api_pbehavior": 15,
                "listalarm_pbehavior": 1,
                "serviceweather_entityComment": 1,
                "serviceweather_entityValidate": 1,
                "serviceweather_addFilter": 1,
                "piechart_interval": 1,
                "api_view": 4,
                "api_map": 2,
                "ea2516be-254c-47ab-9bd5-064584ac743c": 1,
                "alarmStatistics_addFilter": 1,
                "06d65339-74fb-45e1-a961-cd9034594611": 1,
                "userStatistics_addFilter": 1,
                "api_entity": 4,
                "listalarm_listFilters": 1,
                "listalarm_editRemediationInstructionsFilter": 1,
                "crudcontext_createEntity": 1,
                "piechart_sampling": 1,
                "api_declare_ticket_execution": 1,
                "listalarm_listPbehavior": 1,
                "serviceweather_entityInvalidate": 1,
                "listalarm_fastAck": 1,
                "serviceweather_userFilter": 1,
                "listalarm_fastRemoveAlarm": 1,
                "listalarm_groupRequest": 1,
                "listalarm_addBookmark": 1,
                "listalarm_filterByBookmark": 1,
                "listalarm_removeBookmark": 1,
                "crudcontext_delete": 1,
                "map_alarmsList": 1,
                "numbers_userFilter": 1,
                "api_alarm_tag": 4,
                "api_broadcast_message": 4,
                "api_color_theme": 4,
                "api_share_token": 4
            }
        }
    ]);
}
