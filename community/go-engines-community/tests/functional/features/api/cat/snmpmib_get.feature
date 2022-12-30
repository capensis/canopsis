Feature: Get SNMP MIB list
    I need to be able to get a list of MIB
    Admin should have a permission to get a MIB and noperms account shouldn't

    Scenario: given get list request and no auth user shouldn't allow access
        When I do GET /api/v4/cat/snmpmibs
        Then the response code should be 401

    Scenario: given get list request and auth user without permissions shouldn't allow access
        When I am noperms
        When I do GET /api/v4/cat/snmpmibs
        Then the response code should be 403

    Scenario: given get list request should return distinct MIB moduleName
        When I am admin
        When I do GET /api/v4/cat/snmpmibs?nodetype=notification&search=mib-to-get-distinct&projection=moduleName&distinct=true
        Then the response code should be 200
        Then the response body should be:
        """json
        {
            "data": [
                {
                    "moduleName": "mib-to-get-distinct-1997"
                },
                {
                    "moduleName": "mib-to-get-distinct-2002"
                },
                {
                    "moduleName": "mib-to-get-distinct-2012"
                }
            ],
            "meta": {
                "page": 1,
                "per_page": 10,
                "page_count": 1,
                "total_count": 3
            }
        }
        """

    Scenario: given get list request should return distinct MIB moduleName
        When I am admin
        When I do GET /api/v4/cat/snmpmibs?nodetype=notification&moduleName=mib-to-list-module
        Then the response code should be 200
        Then the response body should be:
        """json
        {
            "data": [
                {
                    "status": "current", "moduleName": "mib-to-list-module",
                    "nodetype": "notification", 
                    "description": "mib-to-list-module description", 
                    "oid": "0.1.1", 
                    "objects": 
                        {
                            "nSvcField1": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField2": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField3": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField4": {"nodetype": "object", "module": "mib-to-list-module"}
                        }, 
                    "_id": "0.1.1", "name": "nSvcObject1"
                }, 
                {
                    "status": "current", "moduleName": "mib-to-list-module", 
                    "nodetype": "notification", 
                    "description": "mib-to-list-module description", 
                    "oid": "0.1.2", 
                    "objects": 
                        {
                            "nSvcField5": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField6": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField7": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField8": {"nodetype": "object", "module": "mib-to-list-module"}
                        }, 
                    "_id": "0.1.2", "name": "nSvcObject2"
                }, 
                {
                    "status": "current", "moduleName": "mib-to-list-module", "nodetype": "notification", 
                    "description": "mib-to-list-module description",
                    "oid": "0.1.3", 
                    "objects": 
                        {
                            "nSvcField9": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField10": {"nodetype": "object", "module": "mib-to-list-module"}
                        }, 
                    "_id": "0.1.3", "name": "nSvcObject3"
                }, 
                {
                    "status": "current", "moduleName": "mib-to-list-module", "nodetype": "notification", 
                    "description": "mib-to-list-module description", 
                    "oid": "0.1.4", 
                    "objects": 
                        {
                            "nSvcField11": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField12": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField13": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField14": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField15": {"nodetype": "object", "module": "mib-to-list-module"}, 
                            "nSvcField16": {"nodetype": "object", "module": "mib-to-list-module"}
                        }, 
                        "_id": "0.1.4", "name": "nSvcObject4"
                }
            ],
            "meta": {
                "page": 1,
                "per_page": 10,
                "page_count": 1,
                "total_count": 4
            }
        }
        """
