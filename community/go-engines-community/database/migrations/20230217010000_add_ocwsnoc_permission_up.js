if (!db.default_rights.findOne({_id: "api_ocws_noc"})) {
    db.default_rights.insertOne({
        _id: "api_ocws_noc",
        crecord_name: "api_ocws_noc",
        crecord_type: "action",
        type: "CRUD",
        desc: "OCWS-NOC"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_ocws_noc": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
db.getCollection("ocws_noc_snow_user_new").createIndex({ "sys_id": 1, "provenance": 1 }, {name: "sys_id_1_provenance_1"})
db.getCollection("ocws_noc_snow_location_group_new").createIndex({ "sys_id": 1, "provenance": 1 }, {name: "sys_id_1_provenance_1"})
db.getCollection("ocws_noc_snow_sys_choice_new").createIndex({ "sys_id": 1, "provenance": 1 }, {name: "sys_id_1_provenance_1"})
db.getCollection("default_entities").createIndex({ "type": 1 }, {"name": "type_1"})
