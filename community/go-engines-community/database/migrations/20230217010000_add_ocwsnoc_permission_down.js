db.default_rights.deleteOne({_id: "api_ocws_noc"});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_ocws_noc": "",
    }
});
db.getCollection("ocws_noc_snow_user_new").dropIndex("sys_id_1_provenance_1")
db.getCollection("ocws_noc_snow_location_group_new").dropIndex("sys_id_1_provenance_1")
db.getCollection("ocws_noc_snow_sys_choice_new").dropIndex("sys_id_1_provenance_1")
db.getCollection("default_entities").dropIndex("type_1")
