Feature: Upload a MIB file
  I need to be able to upload a MIB file
  Only admin should be able to upload a MIB file

  Scenario: given upload request and no auth user should not allow access
    When I do POST /api/v4/cat/snmpmibs
    Then the response code should be 401

  Scenario: given upload request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/snmpmibs
    Then the response code should be 403

  Scenario: given upload request with wrong MIB should return Bad Request
    When I am admin
    When I do POST /api/v4/cat/snmpmibs:
    """
    {
      "filecontent": [
        {
          "filename": "MibFile",
          "data": "WRONGMIBDARLING"
        }
      ]
    }
    """
    Then the response code should be 400

  Scenario: given upload request with valid MIB should return OK result
    When I am admin
    When I read file TEST-MIB as testMIB
    When I do POST /api/v4/cat/snmpmibs:
    """
    {
      "filecontent": [
        {
          "filename": "MibFile",
          "data": "{{ .testMIB | json }}"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "counts": {
        "notification": 1,
        "object": 6
      }
    }
    """
    When I do GET /api/v4/cat/snmpmibs?nodetype=node&moduleName=TESTOBJECT-NOTIFY-MIB
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "1.3.6.1.4.1.8072.9999.9999",
          "nodetype": "node",
          "moduleName": "TESTOBJECT-NOTIFY-MIB",
          "oid": "1.3.6.1.4.1.8072.9999.9999",
          "status": "current"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/snmpmibs?nodetype=notification&moduleName=TESTOBJECT-NOTIFY-MIB
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id" : "1.3.6.1.4.1.8072.9999.9999.2",
          "oid" : "1.3.6.1.4.1.8072.9999.9999.2",
          "name" : "tObjectEvent",
          "description" : "The SNMP trap that is generated as a result of an event with object data",
          "objects" : {
            "tObjectname" : {
              "nodetype" : "object",
              "module" : "TESTOBJECT-NOTIFY-MIB"
            },
            "tObjectDurationSec" : {
              "nodetype" : "object",
              "module" : "TESTOBJECT-NOTIFY-MIB"
            }
          },
          "status" : "current",
          "nodetype" : "notification",
          "moduleName" : "TESTOBJECT-NOTIFY-MIB"
        }
      ]
    }
    """
