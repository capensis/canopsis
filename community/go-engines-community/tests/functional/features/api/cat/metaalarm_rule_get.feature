Feature: Get a metaalarm-rule
  I need to be able to get a metaalarm-rule
  Only admin should be able to get a metaalarm-rule

  Scenario: given search request should return metaalarm-rule
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules?search=test-metaalarm-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-metaalarm-to-get-1",
          "auto_resolve": false,
          "config": {
            "time_interval": {
              "value": 10,
              "unit": "s"
            }
          },
          "name": "Test alarm get",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "type": "complex",
          "output_template": "{{ `Rule: {{ .Rule.ID }}; Count: {{ .Count }}; Children: {{ .Children.Alarm.Value.Component }}` }}",
          "created": 1592215337,
          "updated": 1592215337,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-1-pattern"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_total_entity_patterns": null,
          "old_alarm_patterns": null,
          "old_event_patterns": null
        },
        {
          "_id": "test-metaalarm-to-get-2",
          "auto_resolve": false,
          "config": {
            "time_interval": {
              "value": 10,
              "unit": "s"
            }
          },
          "name": "Test alarm get",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "type": "complex",
          "output_template": "{{ `Rule: {{ .Rule.ID }}; Count: {{ .Count }}; Children: {{ .Children.Alarm.Value.Component }}` }}",
          "created": 1592215337,
          "updated": 1592215337,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-2-pattern"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_total_entity_patterns": null,
          "old_alarm_patterns": null,
          "old_event_patterns": null
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

  Scenario: given get request should return metaalarm-rule
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-metaalarm-to-get-1",
      "auto_resolve": false,
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        }
      },
      "name": "Test alarm get",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": "complex",
      "output_template": "{{ `Rule: {{ .Rule.ID }}; Count: {{ .Count }}; Children: {{ .Children.Alarm.Value.Component }}` }}",
      "created": 1592215337,
      "updated": 1592215337,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-get-1-pattern"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_alarm_patterns": null,
      "old_event_patterns": null
    }
    """

  Scenario: given get request should return metaalarm rule with old patterns
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-backward-compatibility-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-metaalarm-rule-backward-compatibility-to-get",
      "auto_resolve": false,
      "config": {
        "threshold_count": 3
      },
      "name": "test-metaalarm-rule-backward-compatibility-to-get-name",
      "type": "complex",
      "old_alarm_patterns": [
        {
          "v": {
            "component": "test-metaalarm-rule-backward-compatibility-component-to-get"
          }
        }
      ],
      "old_entity_patterns": null,
      "old_total_entity_patterns": null,
      "old_event_patterns": null
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/metaalarmrules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metaalarmrules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
