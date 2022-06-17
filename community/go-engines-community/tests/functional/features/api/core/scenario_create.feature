Feature: Create a scenario
  I need to be able to create a scenario
  Only admin should be able to create a scenario

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-1-name",
      "enabled": true,
      "priority": 20,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "author": "test-scenario-to-create-1-action-1-author",
            "output": "test-scenario-to-create-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-2-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-2-resource"
            }
          ],
          "type": "webhook",
          "parameters": {
            "author": "test-scenario-to-create-1-action-2-author",
            "request": {
              "method": "POST",
              "url": "http://test-scenario-to-create-1-action-2-url.com",
              "auth": {
                "username": "test-scenario-to-create-1-action-2-username",
                "password": "test-scenario-to-create-1-action-2-password"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"test-scenario-to-create-1-action-2-payload\": \"test-scenario-to-create-1-action-2-paload-value\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-1-action-2-ticket",
              "test-scenario-to-create-1-action-2-info": "test-scenario-to-create-1-action-2-info-value"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-3-resource"
            }
          ],
          "type": "pbehavior",
          "parameters": {
            "name": "test-scenario-to-create-1-action-3-name",
            "rrule": "FREQ=DAILY",
            "reason": "test-reason-to-edit-scenario",
            "type": "test-type-to-edit-scenario",
            "start_on_trigger": true,
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-create-1-name",
      "author": "root",
      "enabled": true,
      "priority": 20,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "author": "test-scenario-to-create-1-action-1-author",
            "output": "test-scenario-to-create-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-2-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-2-resource"
            }
          ],
          "type": "webhook",
          "parameters": {
            "author": "test-scenario-to-create-1-action-2-author",
            "request": {
              "method": "POST",
              "url": "http://test-scenario-to-create-1-action-2-url.com",
              "auth": {
                "username": "test-scenario-to-create-1-action-2-username",
                "password": "test-scenario-to-create-1-action-2-password"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"test-scenario-to-create-1-action-2-payload\": \"test-scenario-to-create-1-action-2-paload-value\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-1-action-2-ticket",
              "test-scenario-to-create-1-action-2-info": "test-scenario-to-create-1-action-2-info-value"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-3-resource"
            }
          ],
          "type": "pbehavior",
          "parameters": {
            "duration": {
              "value": 3,
              "unit": "s"
            },
            "name": "test-scenario-to-create-1-action-3-name",
            "reason": {
              "_id": "test-reason-to-edit-scenario",
              "description": "test-reason-to-edit-scenario-description",
              "name": "test-reason-to-edit-scenario-name"
            },
            "rrule": "FREQ=DAILY",
            "start_on_trigger": true,
            "type": {
              "_id": "test-type-to-edit-scenario",
              "description": "test-type-to-edit-scenario-description",
              "icon_name": "test-type-to-edit-scenario-icon",
              "name": "test-type-to-edit-scenario-name",
              "priority": 26,
              "type": "maintenance"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-2-name",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test-scenario-to-create-2-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-2-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-2-resource"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "http://test-scenario-to-create-2-action-2-url.com",
              "auth": {
                "username": "test-scenario-to-create-2-action-2-username",
                "password": "test-scenario-to-create-2-action-2-password"
              },
              "headers": {"Content-Type": "application/json"},
              "skip_verify": true,
              "payload": "{\"test-scenario-to-create-2-action-2-payload\": \"test-scenario-to-create-2-action-2-paload-value\"}"
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-2-action-2-ticket",
              "test-scenario-to-create-2-action-2-info": "test-scenario-to-create-2-action-2-info-value"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-3-resource"
            }
          ],
          "type": "pbehavior",
          "parameters": {
            "name": "test-scenario-to-create-2-action-3-name",
            "rrule": "FREQ=DAILY",
            "reason": "test-reason-to-edit-scenario",
            "type": "test-type-to-edit-scenario",
            "start_on_trigger": true,
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/scenarios/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-create-2-name",
      "author": "root",
      "enabled": true,
      "priority": 21,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test-scenario-to-create-2-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-2-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-2-resource"
            }
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "http://test-scenario-to-create-2-action-2-url.com",
              "auth": {
                "username": "test-scenario-to-create-2-action-2-username",
                "password": "test-scenario-to-create-2-action-2-password"
              },
              "headers": {"Content-Type": "application/json"},
              "payload": "{\"test-scenario-to-create-2-action-2-payload\": \"test-scenario-to-create-2-action-2-paload-value\"}",
              "skip_verify": true
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-2-action-2-ticket",
              "test-scenario-to-create-2-action-2-info": "test-scenario-to-create-2-action-2-info-value"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-2-action-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-2-action-3-resource"
            }
          ],
          "type": "pbehavior",
          "parameters": {
            "duration": {
              "value": 3,
              "unit": "s"
            },
            "name": "test-scenario-to-create-2-action-3-name",
            "reason": {
              "_id": "test-reason-to-edit-scenario",
              "description": "test-reason-to-edit-scenario-description",
              "name": "test-reason-to-edit-scenario-name"
            },
            "rrule": "FREQ=DAILY",
            "start_on_trigger": true,
            "type": {
              "_id": "test-type-to-edit-scenario",
              "description": "test-type-to-edit-scenario-description",
              "icon_name": "test-type-to-edit-scenario-icon",
              "name": "test-type-to-edit-scenario-name",
              "priority": 26,
              "type": "maintenance"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """

  Scenario: given create delay request should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-3-name",
      "enabled": true,
      "priority": 12,
      "triggers": ["create"],
      "delay": {
        "value": 3,
        "unit": "s"
      },
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-3-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-create-3-name",
      "author": "root",
      "enabled": true,
      "priority": 12,
      "triggers": ["create"],
      "delay": {
        "value": 3,
        "unit": "s"
      },
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-3-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-3-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/scenarios
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/scenarios
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "priority": 123
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions": "Actions is missing.",
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "triggers": "Triggers is missing."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-check-unique-name-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-4-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-4-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  Scenario: given create request with already exists priority should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 2,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-4-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-4-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "priority": "Priority already exists."
      }
    }
    """

  Scenario: given create request with invalid action should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.alarm_patterns": "AlarmPatterns is missing.",
        "actions.0.drop_scenario_if_not_matched": "DropScenarioIfNotMatched is missing.",
        "actions.0.emit_trigger": "EmitTrigger is missing.",
        "actions.0.entity_patterns": "EntityPatterns is missing.",
        "actions.0.type": "Type is missing."
      }
    }
    """

  Scenario: given create request with action without patterns should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.alarm_patterns": "AlarmPatterns is missing.",
        "actions.0.entity_patterns": "EntityPatterns is missing."
      }
    }
    """

  Scenario: given create request with invalid alarm patterns and empty should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "priority": 21,
      "enabled": true,
      "triggers": [
        "statedec"
      ],
      "disable_during_periods": [],
      "actions": [
        {
          "type": "webhook",
          "parameters": {
            "skip_verify": false,
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false
            },
            "request": {
              "method": "POST",
              "url": "http://localhost:5000",
              "headers": {},
              "payload": "{}"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "alarm_patterns": [
            {
              "component": "component_recette_retry_webhooks"
            }
          ],
          "entity_patterns": []
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "actions.0.alarm_patterns": "Invalid alarm pattern list."
      }
    }
    """

  Scenario: given create request with action with invalid patterns should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [{}],
          "entity_patterns": [{}],
          "type": "snooze",
          "parameters": {
            "output": "test snooze",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.alarm_patterns": "alarm pattern list contains an empty pattern.",
        "actions.0.entity_patterns": "entity pattern list contains an empty pattern."
      }
    }
    """

  Scenario: given create request with snooze action with invalid params should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "snooze",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.parameters.duration": "Duration is missing."
      }
    }
    """

  Scenario: given create request with assocticket action with invalid params should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "assocticket",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.parameters.ticket": "Ticket is missing."
      }
    }
    """

  Scenario: given create request with changestate action with invalid params should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "changestate",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.parameters.state": "State is missing."
      }
    }
    """

  Scenario: given create request with pbehavior action with invalid params should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "pbehavior",
          "parameters": {
            "type": "test-type-not-exist",
            "reason": "test-type-not-exist"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
         "actions.0.parameters.name": "Name is missing.",
         "actions.0.parameters.reason": "Reason doesn't exist.",
         "actions.0.parameters.type": "Type doesn't exist.",
         "actions.0.parameters.start_on_trigger": "StartOnTrigger or Tstart is required.",
         "actions.0.parameters.tstart": "Tstart or StartOnTrigger is required."
      }
    }
    """

  Scenario: given create request with webhook action with invalid params should return error
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "webhook",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.parameters.request": "Request is missing."
      }
    }
    """
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
      "enabled": true,
      "priority": 13,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-resource"
            }
          ],
          "type": "webhook",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "parameters": {
            "request": {}
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "actions.0.parameters.request.method": "Method is missing.",
        "actions.0.parameters.request.url": "URL is missing."
      }
    }
    """

  Scenario: given create request with custom_id should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "my_scenario",
      "name": "my_scenario-name",
      "enabled": true,
      "priority": 987654,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test-scenario-to-create-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "my_scenario"
    }
    """
    When I do GET /api/v4/scenarios/my_scenario
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "my_scenario"
    }
    """

  Scenario: given create request with custom_id should be failed, because of existing id
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-to-check-id",
      "name": "my_scenario-name",
      "enabled": true,
      "priority": 987653,
      "triggers": ["create"],
      "actions": [
        {
          "alarm_patterns": [
            {
              "_id": "test-scenario-to-create-1-action-1-alarm"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-scenario-to-create-1-action-1-resource"
            }
          ],
          "type": "snooze",
          "parameters": {
            "output": "test-scenario-to-create-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with single v.ticket.data pattern should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "ticket-data-scenario-id",
      "name": "webhook scenario",
      "priority": 137,
      "enabled": true,
      "triggers": [
          "create"
      ],
      "disable_during_periods": [],
      "actions": [
          {
            "type": "webhook",
            "parameters": {
                "declare_ticket": {
                    "empty_response": false,
                    "is_regexp": false
                },
                "request": {
                    "method": "POST",
                    "url": "http://localhost:5000",
                    "headers": {},
                    "payload": "{}",
                    "skip_verify": false
                }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false,
            "alarm_patterns": [
                {
                    "v": {
                        "ticket": {
                            "data": {
                                "ticket2_id": {
                                    "regex_match": ".+"
                                }
                            }
                        }
                    }
                }
            ],
            "entity_patterns": []
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "actions": [
          {
              "alarm_patterns": [
                  {
                      "v": {
                          "ticket": {
                              "data": {
                                  "ticket2_id": {
                                      "regex_match": ".+"
                                  }
                              }
                          }
                      }
                  }
              ],
              "entity_patterns": null
          }
      ]
    }
    """
    When I do GET /api/v4/scenarios/ticket-data-scenario-id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "ticket-data-scenario-id"
    }
    """

  Scenario: given create request with single v.parents.is_empty pattern should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "webhook scenario alarm without meta parent",
      "priority": 28,
      "enabled": true,
      "triggers": [
          "create"
      ],
      "disable_during_periods": [],
      "actions": [
          {
            "type": "webhook",
            "parameters": {
                "declare_ticket": {
                    "empty_response": false,
                    "is_regexp": false
                },
                "request": {
                    "method": "POST",
                    "url": "http://localhost:5000",
                    "headers": {},
                    "payload": "{}",
                    "skip_verify": false
                }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false,
            "alarm_patterns": [
                {
                    "_id": "daf712f5-224b-4adc-aeaf-59f37c272fee",
                    "v": {
                        "parents": {
                            "is_empty": true
                        }
                    }
                }
            ],
            "entity_patterns": []
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "actions": [
          {
              "alarm_patterns": [
                  {
                      "_id": "daf712f5-224b-4adc-aeaf-59f37c272fee",
                      "v": {
                          "parents": {
                              "is_empty": true
                          }
                      }
                  }
              ],
              "entity_patterns": null
          }
      ]
    }
    """
    When I do GET /api/v4/scenarios/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "webhook scenario alarm without meta parent"
    }
    """
