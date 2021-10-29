Feature: Get a dynamic-infos
  I need to be able to get a dynamic-infos
  Only admin should be able to get a dynamic-infos

  Scenario: given search request should return dynamic-infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=dynamic_info_
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
          {
              "_id": "dynamic_info_1",
              "entity_patterns": [
                  {
                      "infos": {
                          "alert_name": {
                              "value": {
                                  "regex_match": "(AWS|GCP|HPE)_MAILS_TICKET.*"
                              }
                          }
                      }
                  }
              ],
              "last_modified_date": 1593679995,
              "name": "dynamic_info_1",
              "author": "root",
              "alarm_patterns": null,
              "description": "Consigne pour les alertes reçu depuis un mail des mainteneurs suivants : EMC ECONOCOM ATOS HPSIM NUTANIX PURESTORAGE",
              "enabled": true,
              "disable_during_periods": null,
              "infos": [
                  {
                      "name": "type",
                      "value": "consigne"
                  },
                  {
                      "name": "label",
                      "value": "Consigne sur réception mail Mainteneurs"
                  },
                  {
                      "name": "colibri",
                      "value": "268997"
                  }
              ],
              "creation_date": 1581423405
          },
          {
              "_id": "dynamic_info_2",
              "entity_patterns": [
                  {
                      "infos": {
                          "alert_name": {
                              "value": {
                                  "regex_match": "CENTREONIET_.*IMPORT_POSTES-IET.*"
                              }
                          }
                      }
                  }
              ],
              "last_modified_date": 1593677536,
              "name": "dynamic_info_2",
              "author": "root",
              "alarm_patterns": null,
              "description": "le service d’import des postes dans l’annuaire ne fonctionne pas",
              "enabled": true,
              "disable_during_periods": null,
              "infos": [
                  {
                      "name": "type",
                      "value": "consigne"
                  },
                  {
                      "name": "label",
                      "value": "Consigne Suivi et Administration de l'Infra Neptune"
                  },
                  {
                      "name": "colibri",
                      "value": "148768"
                  }
              ],
              "creation_date": 1581504657
          }
      ],
      "meta": {
          "page": 1,
          "page_count": 1,
          "per_page": 10,
          "total_count": 2
      }
    }
    """

  Scenario: given search DSL request should return dynamic-infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20CENTREONIET
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
          {
              "_id": "dynamic_info_2",
              "entity_patterns": [
                  {
                      "infos": {
                          "alert_name": {
                              "value": {
                                  "regex_match": "CENTREONIET_.*IMPORT_POSTES-IET.*"
                              }
                          }
                      }
                  }
              ],
              "last_modified_date": 1593677536,
              "name": "dynamic_info_2",
              "author": "root",
              "alarm_patterns": null,
              "description": "le service d’import des postes dans l’annuaire ne fonctionne pas",
              "enabled": true,
              "infos": [
                  {
                      "name": "type",
                      "value": "consigne"
                  },
                  {
                      "name": "label",
                      "value": "Consigne Suivi et Administration de l'Infra Neptune"
                  },
                  {
                      "name": "colibri",
                      "value": "148768"
                  }
              ],
              "creation_date": 1581504657,
              "disable_during_periods": null
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

  Scenario: given get request should return dynamic-infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos/dynamic_info_2
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "dynamic_info_2",
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "CENTREONIET_.*IMPORT_POSTES-IET.*"
              }
            }
          }
        }
      ],
      "last_modified_date": 1593677536,
      "name": "dynamic_info_2",
      "author": "root",
      "alarm_patterns": null,
      "description": "le service d’import des postes dans l’annuaire ne fonctionne pas",
      "enabled": true,
      "infos": [
        {
          "name": "type",
          "value": "consigne"
        },
        {
          "name": "label",
          "value": "Consigne Suivi et Administration de l'Infra Neptune"
        },
        {
          "name": "colibri",
          "value": "148768"
        }
      ],
      "creation_date": 1581504657,
      "disable_during_periods": null
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/dynamic-infos
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/dynamic-infos
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """