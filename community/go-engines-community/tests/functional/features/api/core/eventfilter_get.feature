Feature: Get an eventfilter
  I need to be able to get eventfilters
  Only admin should be able to get eventfilters

  Scenario: given search request should return an eventfilter
    When I am admin
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "how it should have ended.",
      "type": "enrichment",
      "old_patterns": null,
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-get-1-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "external_data": {
        "entity": {
          "type": "entity"
        }
      }
    }
    """

  Scenario: given search request should return eventfilters
    When I am admin
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-to-get
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-eventfilter-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "how it should have ended.",
          "type": "enrichment",
          "old_patterns": null,
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-get-1-pattern"
                }
              }
            ]
          ],
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "enabled": true,
          "external_data": {
            "entity": {
              "type": "entity"
            }
          }
        },
        {
          "_id": "test-eventfilter-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "drop filter",
          "type": "drop",
          "old_patterns": null,
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-get-2-pattern"
                }
              }
            ]
          ],
          "config": {},
          "enabled": true
        },
        {
          "_id": "test-eventfilter-to-get-3",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "break filter",
          "type": "break",
          "old_patterns": null,
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-get-3-pattern"
                }
              }
            ]
          ],
          "config": {},
          "enabled": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/eventfilter/rules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/eventfilter/rules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-enrichment-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-enrichment-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given search request should return an eventfilter with old patterns
    When I am admin
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-backward-compatibility-to-get
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-backward-compatibility-to-get",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "how it should have ended.",
      "type": "enrichment",
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "customer",
            "description": "customer",
            "value": "test"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "event_pattern": null,
      "old_patterns": [
        {
          "resource": {
            "regex_match": "test-eventfilter-to-backward-compatibility-to-get"
          }
        }
      ]
    }
    """
