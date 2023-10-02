Feature: Update a SNMP rule
  I need to be able to update a SNMP rule
  Only admin should be able to update a SNMP rule

  Scenario: given update request should update filter
    When I am admin
    Then I do PUT /api/v4/cat/snmprules/test-snmprule-to-update-1:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7000",
            "mibName": "nSvcEvent",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nHostname }}`}}"
        },
        "state": {
            "state": 3,
            "type": "simple"
        },
        "connector_name": {
            "regex": "",
            "formatter": "",
            "value": "test-snmprule-to-update-1-new"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7000",
            "mibName": "nSvcEvent",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "connector_name": {
            "regex": "",
            "formatter": "",
            "value": "test-snmprule-to-update-1-new"
        }
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/cat/snmprules/test-snmprule-to-update-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/snmprules/test-snmprule-to-update-1
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/snmprules/test-snmprule-to-update-not-found:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "resource-1"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.10000.1.1",
            "mibName": "mib1",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "regex": "",
            "formatter": "",
            "value": "component1"
        },
        "state": {
            "state": 1,
            "type": "simple"
        },
        "connector_name": {
            "regex": "",
            "formatter": "",
            "value": "test-snmprule-to-update-not-found"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "output1"
        }
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/snmprules/test-snmprule-to-update-1:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "component.value": "Value of component is missing.",
        "connector_name.value": "Value of connector_name is missing.",
        "state.type": "State type is missing.",
        "oid.oid": "oid is missing.",
        "oid.mibName": "mibName is missing.",
        "oid.moduleName": "moduleName is missing."
      }
    }
    """
