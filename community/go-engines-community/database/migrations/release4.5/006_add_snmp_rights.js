db.default_rights.insertMany([
    {
        "_id": "api_snmprule",
        "loader_id": "api_snmprule",
        "crecord_type": "action",
        "crecord_name": "api_snmprule",
        "desc": "SNMP api",
        "type": "CRUD"
    },
    {
        "_id": "api_snmpmib",
        "loader_id": "api_snmpmib",
        "crecord_type": "action",
        "crecord_name": "api_snmpmib",
        "type": "CRUD",
        "desc": "SNMP MIB api"
    },
]);

db.default_rights.update(
    {
        crecord_name: "admin",
        crecord_type: "role",
    },
    {
        $set: {
            "rights.api_snmprule": {
                checksum: 15,
                crecord_type: "right",
            },
            "rights.api_snmpmib": {
                checksum: 15,
                crecord_type: "right",
            },
        },
    }
);
