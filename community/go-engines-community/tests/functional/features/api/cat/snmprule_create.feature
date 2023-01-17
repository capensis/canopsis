Feature: Create a SNMP rule
  I need to be able to create a SNMP rule
  Only admin should be able to create a SNMP rule

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/snmprules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/snmprules
    Then the response code should be 403

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7",
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
            "value": "test-snmprule-to-create-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7",
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
            "value": "test-snmprule-to-create-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    When I do GET /api/v4/cat/snmprules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7",
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
            "value": "test-snmprule-to-create-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """

  Scenario: given create a SNMP rule request should return validation error
    When I am admin
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
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
            "value": "test-snmprule-to-create-invalid-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "oid.oid": "oid is missing.",
        "oid.mibName": "mibName is missing.",
        "oid.moduleName": "moduleName is missing."
      }
    }
    """
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
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
            "value": "test-snmprule-to-create-invalid-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        },
        "oid": {
            "mibName": "nSvcEvent",
            "moduleName": "TEST-NOTIFY-MIB",
            "oid": "test-snmprule-to-create-invalid-connector-test-snmprule-to-create-invalid-connector-test-snmprule-to-create-invalid-connector-test-snmprule-to-create-invalid-connector-test-snmprule-to-create-invalid-connector-test-snmprule-to-create-invalid-connector-123456789"
        }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "oid.oid": "oid should be 255 or less."
      }
    }
    """
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7",
            "mibName": "nSvcEvent",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nHostname }}`}}"
        },
        "connector_name": {
            "regex": "",
            "formatter": "",
            "value": "test-snmprule-to-create-invalid-connector"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "state.type": "State type is missing."
      }
    }
    """
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "resource": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcDesc }}`}}"
        },
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.7",
            "mibName": "nSvcEvent",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "connector_name": {
            "regex": "",
            "formatter": "",
            "value": "test-snmprule-to-create-invalid-connector"
        },
        "state": {
            "state": 3,
            "type": "simple"
        },
        "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
        }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "component.value": "Value of component is missing."
      }
    }
    """

  Scenario: given create request with template state and stateoid should return ok
    When I am admin
    When I do POST /api/v4/cat/snmprules:
    """json
    {
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.6",
            "mibName": "nHostNotify",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "value": "prefix-1{{`{{ nHostAckAuthor }}`}}",
            "regex": "rx",
            "formatter": "gdsf"
        },
        "connector_name": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "output": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "resource": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "state": {
            "type": "template",
            "state": 0,
            "stateoid": {
                "value": "prefix2{{`{{ nHostAckAuthor }}`}}",
                "regex": "",
                "formatter": ""
            },
            "info": "fsad",
            "minor":"fds",
            "major":"asd",
            "critical":"asd"
        }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.6",
            "mibName": "nHostNotify",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "value": "prefix-1{{`{{ nHostAckAuthor }}`}}",
            "regex": "rx",
            "formatter": "gdsf"
        },
        "connector_name": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "output": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "resource": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "state": {
            "type": "template",
            "state": 0,
            "stateoid": {
                "value": "prefix2{{`{{ nHostAckAuthor }}`}}",
                "regex": "",
                "formatter": ""
            },
            "info": "fsad",
            "minor":"fds",
            "major":"asd",
            "critical":"asd"
        }
    }
    """
    When I do GET /api/v4/cat/snmprules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "oid": {
            "oid": "1.3.6.1.4.1.20006.1.6",
            "mibName": "nHostNotify",
            "moduleName": "NAGIOS-NOTIFY-MIB"
        },
        "component": {
            "value": "prefix-1{{`{{ nHostAckAuthor }}`}}",
            "regex": "rx",
            "formatter": "gdsf"
        },
        "connector_name": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "output": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "resource": {
            "value": "{{`{{ nHostAckComment }}`}}",
            "regex": "",
            "formatter": ""
        },
        "state": {
            "type": "template",
            "state": 0,
            "stateoid": {
                "value": "prefix2{{`{{ nHostAckAuthor }}`}}",
                "regex": "",
                "formatter": ""
            },
            "info": "fsad",
            "minor":"fds",
            "major":"asd",
            "critical":"asd"
        }
    }
    """
