Feature: Get a SNMP rule
    I need to be able to get a SNMP rule
    Admin should have a permission to get a SNMP rule and noperms account shouldn't

    Scenario: given get list request and no auth user shouldn't allow access
        When I do GET /api/v4/cat/snmprules
        Then the response code should be 401

    Scenario: given get list request and auth user without permissions shouldn't allow access
        When I am noperms
        When I do GET /api/v4/cat/snmprules
        Then the response code should be 403

    Scenario: given get list request should return SNMP rules
        When I am admin
        When I do GET /api/v4/cat/snmprules?search=test-snmprule-to-get
        Then the response code should be 200
        Then the response body should be:
        """json
        {
          "data": [
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
                "value": "test-snmprule-to-get-connector"
              },
              "output": {
                "regex": "",
                "formatter": "",
                "value": "{{`{{ nSvcOutput }}`}}"
              },
              "_id": "test-snmprule-to-get"
            }
          ],
          "meta": {
            "page": 1,
            "page_count": 1,
            "per_page": 10,
            "total_count": 1
          }
        }
        """

    Scenario: given get request by ID should return SNMP rule
        When I am admin
        When I do GET /api/v4/cat/snmprules/test-snmprule-to-get
        Then the response code should be 200
        Then the response body should be:
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
            "value": "test-snmprule-to-get-connector"
          },
          "output": {
            "regex": "",
            "formatter": "",
            "value": "{{`{{ nSvcOutput }}`}}"
          },
          "_id": "test-snmprule-to-get"
        }
        """
        When I am admin
        When I do GET /api/v4/cat/snmprules/test-snmprule-not-found
        Then the response code should be 404
        Then the response body should be:
        """json
        {
          "error": "Not found"
        }
        """