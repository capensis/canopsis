Feature: Get a scenario
  I need to be able to read a scenario
  Only admin should be able to read a scenario

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/scenarios
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/scenarios
    Then the response code should be 403

  Scenario: given get all request should return scenarios
    When I am admin
    When I do GET /api/v4/scenarios?search=test-scenario-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-scenario-to-get-1",
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-get-1-action-1-alarm"
                    }
                  }
                ]
              ],
              "old_entity_patterns": null,
              "old_alarm_patterns": null,
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false,
              "type": "ack",
              "parameters": {
                "author": "test-scenario-to-get-1-action-1-author",
                "output": "test-scenario-to-get-1-action-1-output"
              },
              "comment": "test-scenario-to-get-1-action-1-comment"
            },
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-get-1-action-2-alarm"
                    }
                  }
                ]
              ],
              "old_entity_patterns": null,
              "old_alarm_patterns": null,
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false,
              "type": "pbehavior",
              "parameters": {
                "name": "test-scenario-to-get-1-action-2-name",
                "author": "test-scenario-to-get-1-action-2-author",
                "rrule": "FREQ=DAILY",
                "reason": {
                  "_id": "test-reason-to-edit-scenario",
                  "description": "test-reason-to-edit-scenario-description",
                  "name": "test-reason-to-edit-scenario-name",
                  "created": 1592215337
                },
                "type": {
                  "_id": "test-type-to-edit-scenario",
                  "description": "test-type-to-edit-scenario-description",
                  "icon_name": "test-type-to-edit-scenario-icon",
                  "name": "test-type-to-edit-scenario-name",
                  "priority": 25,
                  "type": "maintenance"
                },
                "start_on_trigger": true,
                "duration": {
                  "value": 3,
                  "unit": "s"
                }
              },
              "comment": "test-scenario-to-get-1-action-2-comment"
            }
          ],
          "enabled": true,
          "name": "test-scenario-to-get-1-name",
          "triggers": [
            "create"
          ]
        },
        {
          "_id": "test-scenario-to-get-2",
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-get-2-action-1-alarm"
                    }
                  }
                ]
              ],
              "old_entity_patterns": null,
              "old_alarm_patterns": null,
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false,
              "type": "ack",
              "parameters": {
                "author": "test-scenario-to-get-2-action-1-author",
                "output": "test-scenario-to-get-2-action-1-output"
              },
              "comment": ""
            }
          ],
          "delay": {
            "value": 10,
            "unit": "s"
          },
          "disable_during_periods": null,
          "enabled": true,
          "name": "test-scenario-to-get-2-name",
          "triggers": [
            "create"
          ]
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

  Scenario: given sort request should return sorted scenarios
    When I am admin
    When I do GET /api/v4/scenarios?search=test-scenario-to-get&sort_by=name&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-scenario-to-get-2"
        },
        {
          "_id": "test-scenario-to-get-1"
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

  Scenario: given get request should return scenario
    When I am admin
    When I do GET /api/v4/scenarios/test-scenario-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-to-get-1",
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-get-1-action-1-alarm"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_alarm_patterns": null,
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "type": "ack",
          "parameters": {
            "author": "test-scenario-to-get-1-action-1-author",
            "output": "test-scenario-to-get-1-action-1-output"
          },
          "comment": "test-scenario-to-get-1-action-1-comment"
        },
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-get-1-action-2-alarm"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_alarm_patterns": null,
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "type": "pbehavior",
          "parameters": {
            "name": "test-scenario-to-get-1-action-2-name",
            "author": "test-scenario-to-get-1-action-2-author",
            "rrule": "FREQ=DAILY",
            "reason": {
              "_id": "test-reason-to-edit-scenario",
              "description": "test-reason-to-edit-scenario-description",
              "name": "test-reason-to-edit-scenario-name",
              "created": 1592215337
            },
            "type": {
              "_id": "test-type-to-edit-scenario",
              "description": "test-type-to-edit-scenario-description",
              "icon_name": "test-type-to-edit-scenario-icon",
              "name": "test-type-to-edit-scenario-name",
              "priority": 25,
              "type": "maintenance"
            },
            "start_on_trigger": true,
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "comment": "test-scenario-to-get-1-action-2-comment"
        }
      ],
      "enabled": true,
      "name": "test-scenario-to-get-1-name",
      "triggers": [
        "create"
      ]
    }
    """

  Scenario: given get request should return scenario with old patterns
    When I am admin
    When I do GET /api/v4/scenarios/test-scenario-backward-compatibility-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-backward-compatibility-to-get-1",
      "name": "test-scenario-backward-compatibility-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "disable_during_periods": null,
      "triggers": [
          "create"
      ],
      "actions": [
        {
          "type": "ack",
          "comment": "",
          "parameters": {
            "output": "test-scenario-backward-compatibility-to-get-1-action-1-output",
            "author": "test-scenario-backward-compatibility-to-get-1-action-1-author"
          },
          "old_alarm_patterns": [
            {
              "_id": "test-scenario-backward-compatibility-to-get-1-action-1-alarm"
            }
          ],
          "old_entity_patterns": [
            {
              "name": "test-scenario-backward-compatibility-to-get-1-action-1-name"
            }
          ],
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/scenarios/notexist
    Then the response code should be 404
