Feature: Create a scenario
  I need to be able to bulk create scenarios
  Only admin should be able to bulk create scenarios

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/bulk/scenarios
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/scenarios
    Then the response code should be 403

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/bulk/scenarios:
    """
    [
      {
        "name": "test-scenario-to-bulk-create-1-name",
        "enabled": true,
        "priority": 10,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-action-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-action-1-resource"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test-scenario-to-bulk-create-1-action-1-output",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          },
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-action-2-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-action-2-resource"
              }
            ],
            "type": "webhook",
            "parameters": {
              "request": {
                "method": "POST",
                "url": "http://test-scenario-to-bulk-create-1-action-2-url.com",
                "auth": {
                  "username": "test-scenario-to-bulk-create-1-action-2-username",
                  "password": "test-scenario-to-bulk-create-1-action-2-password"
                },
                "headers": {"Content-Type": "application/json"},
                "payload": "{\"test-scenario-to-bulk-create-1-action-2-payload\": \"test-scenario-to-bulk-create-1-action-2-paload-value\"}"
              },
              "declare_ticket": {
                "empty_response": false,
                "is_regexp": false,
                "ticket_id": "test-scenario-to-bulk-create-1-action-2-ticket",
                "test-scenario-to-bulk-create-1-action-2-info": "test-scenario-to-bulk-create-1-action-2-info-value"
              },
              "retry_count": 3,
              "retry_delay": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          },
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-action-3-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-action-3-resource"
              }
            ],
            "type": "pbehavior",
            "parameters": {
              "name": "test-scenario-to-bulk-create-1-action-3-name",
              "rrule": "FREQ=DAILY",
              "reason": "test-reason-to-edit-scenario",
              "type": "test-type-to-edit-scenario",
              "start_on_trigger": true,
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
        "name": "test-scenario-to-bulk-create-2-name",
        "enabled": true,
        "priority": 11,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-2-action-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-2-action-1-resource"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test-scenario-to-bulk-create-2-action-1-output",
              "duration": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          },
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-2-action-2-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-2-action-2-resource"
              }
            ],
            "type": "webhook",
            "parameters": {
              "request": {
                "method": "POST",
                "url": "http://test-scenario-to-bulk-create-2-action-2-url.com",
                "auth": {
                  "username": "test-scenario-to-bulk-create-2-action-2-username",
                  "password": "test-scenario-to-bulk-create-2-action-2-password"
                },
                "headers": {"Content-Type": "application/json"},
                "skip_verify": true,
                "payload": "{\"test-scenario-to-bulk-create-2-action-2-payload\": \"test-scenario-to-bulk-create-2-action-2-paload-value\"}"
              },
              "declare_ticket": {
                "empty_response": false,
                "is_regexp": false,
                "ticket_id": "test-scenario-to-bulk-create-2-action-2-ticket",
                "test-scenario-to-bulk-create-2-action-2-info": "test-scenario-to-bulk-create-2-action-2-info-value"
              },
              "retry_count": 3,
              "retry_delay": {
                "seconds": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          },
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-2-action-3-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-2-action-3-resource"
              }
            ],
            "type": "pbehavior",
            "parameters": {
              "name": "test-scenario-to-bulk-create-2-action-3-name",
              "rrule": "FREQ=DAILY",
              "reason": "test-reason-to-edit-scenario",
              "type": "test-type-to-edit-scenario",
              "start_on_trigger": true,
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
    When I do GET /api/v4/scenarios?search=test-scenario-to-bulk-create&sort=asc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-scenario-to-bulk-create-1-name",
          "author": "root",
          "enabled": true,
          "priority": 10,
          "triggers": ["create"],
          "actions": [
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-1-action-1-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-1-action-1-resource"
                }
              ],
              "type": "snooze",
              "parameters": {
                "author": "root",
                "user": "root",
                "output": "test-scenario-to-bulk-create-1-action-1-output",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            },
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-1-action-2-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-1-action-2-resource"
                }
              ],
              "type": "webhook",
              "parameters": {
                "request": {
                  "method": "POST",
                  "url": "http://test-scenario-to-bulk-create-1-action-2-url.com",
                  "auth": {
                    "username": "test-scenario-to-bulk-create-1-action-2-username",
                    "password": "test-scenario-to-bulk-create-1-action-2-password"
                  },
                  "headers": {"Content-Type": "application/json"},
                  "payload": "{\"test-scenario-to-bulk-create-1-action-2-payload\": \"test-scenario-to-bulk-create-1-action-2-paload-value\"}"
                },
                "declare_ticket": {
                  "empty_response": false,
                  "is_regexp": false,
                  "ticket_id": "test-scenario-to-bulk-create-1-action-2-ticket",
                  "test-scenario-to-bulk-create-1-action-2-info": "test-scenario-to-bulk-create-1-action-2-info-value"
                },
                "retry_count": 3,
                "retry_delay": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            },
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-1-action-3-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-1-action-3-resource"
                }
              ],
              "type": "pbehavior",
              "parameters": {
                "author": "root",
                "user": "root",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                },
                "name": "test-scenario-to-bulk-create-1-action-3-name",
                "reason": {
                  "_id": "test-reason-to-edit-scenario",
                  "description": "test-reason-to-edit-scenario-description",
                  "name": "test-reason-to-edit-scenario-name"
                },
                "rrule": "FREQ=DAILY",
                "start_on_trigger": true,
                "tstart": null,
                "tstop": null,
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
        },
        {
          "name": "test-scenario-to-bulk-create-2-name",
          "author": "root",
          "enabled": true,
          "priority": 11,
          "triggers": ["create"],
          "actions": [
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-2-action-1-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-2-action-1-resource"
                }
              ],
              "type": "snooze",
              "parameters": {
                "author": "root",
                "user": "root",
                "output": "test-scenario-to-bulk-create-2-action-1-output",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            },
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-2-action-2-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-2-action-2-resource"
                }
              ],
              "type": "webhook",
              "parameters": {
                "request": {
                  "method": "POST",
                  "url": "http://test-scenario-to-bulk-create-2-action-2-url.com",
                  "auth": {
                    "username": "test-scenario-to-bulk-create-2-action-2-username",
                    "password": "test-scenario-to-bulk-create-2-action-2-password"
                  },
                  "headers": {"Content-Type": "application/json"},
                  "payload": "{\"test-scenario-to-bulk-create-2-action-2-payload\": \"test-scenario-to-bulk-create-2-action-2-paload-value\"}",
                  "skip_verify": true
                },
                "declare_ticket": {
                  "empty_response": false,
                  "is_regexp": false,
                  "ticket_id": "test-scenario-to-bulk-create-2-action-2-ticket",
                  "test-scenario-to-bulk-create-2-action-2-info": "test-scenario-to-bulk-create-2-action-2-info-value"
                },
                "retry_count": 3,
                "retry_delay": {
                  "seconds": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            },
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-create-2-action-3-alarm"
                }
              ],
              "entity_patterns": [
                {
                  "name": "test-scenario-to-bulk-create-2-action-3-resource"
                }
              ],
              "type": "pbehavior",
              "parameters": {
                "author": "root",
                "user": "root",
                "duration": {
                  "seconds": 3,
                  "unit": "s"
                },
                "name": "test-scenario-to-bulk-create-2-action-3-name",
                "reason": {
                  "_id": "test-reason-to-edit-scenario",
                  "description": "test-reason-to-edit-scenario-description",
                  "name": "test-reason-to-edit-scenario-name"
                },
                "rrule": "FREQ=DAILY",
                "start_on_trigger": true,
                "tstart": null,
                "tstop": null,
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/scenarios:
    """
    {
      "_id": "my_scenario_bulk_create",
      "name": "my_scenario_bulk_create-name",
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
              "seconds": 3,
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
    When I do POST /api/v4/bulk/scenarios:
    """
    [
      {},
      {
        "name": "test-scenario-to-check-unique-name-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-4-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-4-resource"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze",
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
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 2,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-4-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-4-resource"
              }
            ],
            "type": "snooze",
            "parameters": {
              "output": "test snooze",
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
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
          }
        ]
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "type": "snooze",
            "parameters": {
              "output": "test snooze",
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
        "name": "test-scenario-to-bulk-create-4-name",
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
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
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
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-resource"
              }
            ],
            "type": "snooze",
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-resource"
              }
            ],
            "type": "assocticket",
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-resource"
              }
            ],
            "type": "changestate",
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-resource"
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
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "priority": 13,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_patterns": [
              {
                "_id": "test-scenario-to-bulk-create-1-alarm"
              }
            ],
            "entity_patterns": [
              {
                "name": "test-scenario-to-bulk-create-1-resource"
              }
            ],
            "type": "webhook",
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
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
        "0.actions": "Actions is missing.",
        "0.enabled": "Enabled is missing.",
        "0.name": "Name is missing.",
        "0.priority": "Priority is missing.",
        "0.triggers": "Triggers is missing.",
        "1.name": "Name already exists.",
        "2.priority": "Priority already exists.",
        "3.actions.0.alarm_patterns": "AlarmPatterns is missing.",
        "3.actions.0.drop_scenario_if_not_matched": "DropScenarioIfNotMatched is missing.",
        "3.actions.0.emit_trigger": "EmitTrigger is missing.",
        "3.actions.0.entity_patterns": "EntityPatterns is missing.",
        "3.actions.0.type": "Type is missing.",
        "4.actions.0.alarm_patterns": "AlarmPatterns is missing.",
        "4.actions.0.entity_patterns": "EntityPatterns is missing.",
        "5.actions.0.alarm_patterns": "Invalid alarm pattern list.",
        "6.actions.0.alarm_patterns": "alarm pattern list contains an empty pattern.",
        "6.actions.0.entity_patterns": "entity pattern list contains an empty pattern.",
        "7.actions.0.parameters.duration": "Duration is missing.",
        "8.actions.0.parameters.ticket": "Ticket is missing.",
        "9.actions.0.parameters.output": "Output is missing.",
        "9.actions.0.parameters.state": "State is missing.",
        "10.actions.0.parameters.name": "Name is missing.",
        "10.actions.0.parameters.reason": "Reason doesn't exist.",
        "10.actions.0.parameters.type": "Type doesn't exist.",
        "10.actions.0.parameters.start_on_trigger": "StartOnTrigger or Tstart is required.",
        "10.actions.0.parameters.tstart": "Tstart or StartOnTrigger is required.",
        "11.actions.0.parameters.request.method": "Method is missing.",
        "11.actions.0.parameters.request.url": "URL is missing.",
        "12._id": "ID already exists."
      }
    }
    """
