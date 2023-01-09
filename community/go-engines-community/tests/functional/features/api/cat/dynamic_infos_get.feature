Feature: Get a dynamic infos
  I need to be able to get a dynamic infos
  Only admin should be able to get a dynamic infos

  Scenario: given search request should return dynamic infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=test-dynamic-infos-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-dynamic-infos-to-get-1",
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-dynamic-infos-to-get-1-alarm-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-dynamic-infos-to-get-1-entity-pattern"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1581423405,
          "description": "test-dynamic-infos-to-get-1-description",
          "disable_during_periods": null,
          "enabled": true,
          "infos": [
            {
              "name": "test-dynamic-infos-to-get-1-info-1-name",
              "value": "test-dynamic-infos-to-get-1-info-1-value"
            },
            {
              "name": "test-dynamic-infos-to-get-1-info-2-name",
              "value": "test-dynamic-infos-to-get-1-info-2-value"
            }
          ],
          "updated": 1593679995,
          "name": "test-dynamic-infos-to-get-1-name"
        },
        {
          "_id": "test-dynamic-infos-to-get-2",
          "old_alarm_patterns": [
            {
              "v": {
                "connector": "test-dynamic-infos-to-get-2-alarm-pattern"
              }
            }
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1581423405,
          "description": "test-dynamic-infos-to-get-2-description",
          "disable_during_periods": null,
          "enabled": true,
          "old_entity_patterns": [
            {
              "_id": "test-dynamic-infos-to-get-2-entity-pattern"
            }
          ],
          "infos": [
            {
              "name": "test-dynamic-infos-to-get-2-info-1-name",
              "value": "test-dynamic-infos-to-get-2-info-1-value"
            },
            {
              "name": "test-dynamic-infos-to-get-2-info-2-name",
              "value": "test-dynamic-infos-to-get-2-info-2-value"
            }
          ],
          "updated": 1593679995,
          "name": "test-dynamic-infos-to-get-2-name"
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

  Scenario: given pattern search request should return dynamic infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"test-dynamic-infos-to-get-2"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {"_id": "test-dynamic-infos-to-get-2"}
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given get request should return dynamic infos
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-dynamic-infos-to-get-2",
      "old_alarm_patterns": [
        {
          "v": {
            "connector": "test-dynamic-infos-to-get-2-alarm-pattern"
          }
        }
      ],
      "old_entity_patterns": [
        {
          "_id": "test-dynamic-infos-to-get-2-entity-pattern"
        }
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1581423405,
      "description": "test-dynamic-infos-to-get-2-description",
      "disable_during_periods": null,
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-get-2-info-1-name",
          "value": "test-dynamic-infos-to-get-2-info-1-value"
        },
        {
          "name": "test-dynamic-infos-to-get-2-info-2-name",
          "value": "test-dynamic-infos-to-get-2-info-2-value"
        }
      ],
      "updated": 1593679995,
      "name": "test-dynamic-infos-to-get-2-name"
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/dynamic-infos
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/dynamic-infos
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
