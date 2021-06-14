Feature: Create an dynamic-infos
  I need to be able to create a dynamic-infos
  Only admin should be able to create a dynamic-infos

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id": "dynamic_infos_3",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NAGIOSPTSR_.*FILESYSTEMS.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_infos_3",
      "alarm_patterns": null,
      "description": "test dynamic infos",
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Proc_NagiosPtsr_Filesystems"
        },
        {
          "name": "colibri",
          "value": "153028"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "dynamic_infos_3",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NAGIOSPTSR_.*FILESYSTEMS.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_infos_3",
      "author": "root",
      "alarm_patterns": null,
      "description": "test dynamic infos",
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Proc_NagiosPtsr_Filesystems"
        },
        {
          "name": "colibri",
          "value": "153028"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/dynamic-infos/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "dynamic_infos_3",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "NAGIOSPTSR_.*FILESYSTEMS.*"
              }
            }
          }
        }
      ],
      "name": "dynamic_infos_3",
      "author": "root",
      "alarm_patterns": null,
      "description": "test dynamic infos",
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Proc_NagiosPtsr_Filesystems"
        },
        {
          "name": "colibri",
          "value": "153028"
        }
      ]
    }
    """

  Scenario: given search DSL request should return dynamic-infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20NAGIOSPTSR
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "dynamic_infos_3",
          "entity_patterns": [
            {
              "infos": {
                "alert_name": {
                  "value": {
                    "regex_match": "NAGIOSPTSR_.*FILESYSTEMS.*"
                  }
                }
              }
            }
          ],
          "name": "dynamic_infos_3",
          "author": "root",
          "alarm_patterns": null,
          "description": "test dynamic infos",
          "infos": [
            {
              "name": "type",
              "value": "consigne"
            },
            {
              "name": "label",
              "value": "Proc_NagiosPtsr_Filesystems"
            },
            {
              "name": "colibri",
              "value": "153028"
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

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id": "test_insert1",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "(ATOS|EMC|ECONOCOM|HPSIM|PURESTORAGE|NUTANIX|HPE)_MAILS_TICKET.*"
              }
            }
          }
        }
      ],
      "name": "Consigne sur réception mail Mainteneur",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "infos": [
        {
          "name": "name"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "infos.0.value": "Value is missing."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id": "test_insert1",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "(ATOS|EMC|ECONOCOM|HPSIM|PURESTORAGE|NUTANIX|HPE)_MAILS_TICKET.*"
              }
            }
          }
        }
      ],
      "last_modified_date": 1593679995,
      "name": "Consigne sur réception mail Mainteneur",
      "author": "root",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "infos": [
        {
          "name": "name",
          "value": "test"
        }
      ],
      "creation_date": 1581423405
    }
    """
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id": "test_insert1",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "(ATOS|EMC|ECONOCOM|HPSIM|PURESTORAGE|NUTANIX|HPE)_MAILS_TICKET.*"
              }
            }
          }
        }
      ],
      "last_modified_date": 1593679995,
      "name": "Consigne sur réception mail Mainteneur",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "infos": [
        {
          "name": "name",
          "value": "test"
        }
      ],
      "creation_date": 1581423405
    }
    """
    Then the response code should be 403

  Scenario: given create request with already exists _id should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id": "dynamic_infos_3",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "(ATOS|EMC|ECONOCOM|HPSIM|PURESTORAGE|NUTANIX|HPE)_MAILS_TICKET.*"
              }
            }
          }
        }
      ],
      "name": "Consigne sur réception mail Mainteneur",
      "alarm_patterns": null,
      "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
      "infos": [
        {
          "name": "name",
          "value": "test"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with wrong alarm pattern should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """
    {
      "_id" : "2a531090-8288-4db5-84f0-32b9d62fd8b2",
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