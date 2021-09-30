Feature: Update a dynamic-infos
  I need to be able to update a dynamic-infos
  Only admin should be able to update a dynamic-infos

  Scenario: given update request should update dynamic-infos
    When I am admin
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update:
    """
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NEW_UPDATE.*"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-update",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne sur réception mail Mainteneurs"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-dynamic-infos-to-update",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NEW_UPDATE.*"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-update",
      "author": "root",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne sur réception mail Mainteneurs"
        }
      ]
    }
    """

  Scenario: given search DSL request should return dynamic-infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20NEW_UPDATE
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-dynamic-infos-to-update",
          "entity_patterns": [
            {
              "infos": {
                "alert_name": {
                  "value": {
                    "regex_match": "NEW_UPDATE.*"
                  }
                }
              }
            }
          ],
          "name": "test-dynamic-infos-to-update",
          "author": "root",
          "alarm_patterns": null,
          "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
          "enabled": true,
          "infos": [
            {
              "name": "type",
              "value": "consigne"
            },
            {
              "name": "label",
              "value": "Consigne sur réception mail Mainteneurs"
            }
          ]
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

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/dynamic-infos/test-alarm-get-metaalarm-rule-1:
    """
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NEW_UPDATE.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_info_1",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne sur réception mail Mainteneurs"
        }
      ],
      "pattern": [
        "(AWS|GCP|HPE)_MAILS_TICKET.*"
      ]
    }
    """
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/dynamic-infos/test-alarm-get-metaalarm-rule-1:
    """
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NEW_UPDATE.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_info_1",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne sur réception mail Mainteneurs"
        }
      ],
      "pattern": [
        "(AWS|GCP|HPE)_MAILS_TICKET.*"
      ]
    }
    """
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found:
    """
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NEW_UPDATE.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_info_1",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne sur réception mail Mainteneurs"
        }
      ]
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with wrong alarm pattern should return error
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update:
    """
    {
      "_id" : "test-dynamic-infos-to-update",
      "entity_patterns" : null,
      "disable_during_periods" : [ ],
      "name" : "Name",
      "alarm_patterns" : [
              {
                      "v" : {

                      }
              }
      ],
      "description" : "Desc",
      "enabled": true,
      "infos" : [
              {
                      "name" : "Name",
                      "value" : ""
              }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "alarm_patterns": "Invalid alarm patterns."
      }
    }
    """