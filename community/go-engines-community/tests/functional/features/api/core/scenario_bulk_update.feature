Feature: Update a scenario
  I need to be able to bulk update scenarios
  Only admin should be able to bulk update scenarios

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/bulk/scenarios
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/scenarios
    Then the response code should be 403


  Scenario: given update request should update exception
    When I am admin
    When I do PUT /api/v4/bulk/scenarios:
    """
    [
      {
        "_id": "test-scenario-to-bulk-update-1",
        "name": "test-scenario-to-bulk-update-1-name",
        "enabled": true,
        "priority": 15,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-update-1-alarm-updated"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze updated",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "_id": "test-scenario-to-bulk-update-2",
        "name": "test-scenario-to-bulk-update-2-name",
        "enabled": true,
        "priority": 15,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-update-2-alarm-updated"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze updated",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/scenarios?search=test-scenario-to-bulk-update&sort=asc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-scenario-to-bulk-update-1",
          "name": "test-scenario-to-bulk-update-1-name",
          "author": "root",
          "enabled": true,
          "priority": 15,
          "delay": null,
          "disable_during_periods": null,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-update-1-alarm-updated"
                }
              ],
              "entity_patterns": null,
              "type": "snooze",
              "parameters": {
                "author": "root",
                "user": "root",
                "output": "test snooze updated",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ]
        },
        {
          "_id": "test-scenario-to-bulk-update-2",
          "name": "test-scenario-to-bulk-update-2-name",
          "author": "root",
          "enabled": true,
          "priority": 15,
          "delay": null,
          "disable_during_periods": null,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-update2-alarm-updated"
                }
              ],
              "entity_patterns": null,
              "type": "snooze",
              "parameters": {
                "author": "root",
                "user": "root",
                "output": "test snooze updated",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
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

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """
    [
      {
        "_id": "test-scenario-to-update-1",
        "name": "test-scenario-to-check-unique-name-name",
        "enabled": true,
        "priority": 16,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-update-1-alarm-updated"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze updated",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "_id": "test-scenario-to-update-1",
        "name": "test-scenario-to-update-1-name",
        "enabled": true,
        "priority": 2,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-update-1-alarm-updated"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze updated",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "0.name": "Name already exists.",
        "1.priority": "Priority already exists."
      }
    }
    """
