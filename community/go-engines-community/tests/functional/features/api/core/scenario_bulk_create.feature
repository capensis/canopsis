Feature: Bulk create scenarios
  I need to be able to bulk create scenarios
  Only admin should be able to bulk create scenarios

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/scenarios
    Then the response code should be 401

  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/scenarios
    Then the response code should be 403

  Scenario: given bulk create request should return multi status and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/scenarios:
    """json
    [
      {
        "_id": "bulk-create-scenario-1",
        "name": "test-scenario-to-bulk-create-1-name",
        "enabled": true,
        "priority": 200000,
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
                "value": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
      {
        "_id": "bulk-create-scenario-1",
        "name": "test-scenario-to-bulk-create-1-name",
        "enabled": true,
        "priority": 200000,
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
                "value": 3,
                "unit": "s"
              }
            },
            "drop_scenario_if_not_matched": false,
            "emit_trigger": false
          }
        ]
      },
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
      },
      [],
      {
        "_id": "bulk-create-scenario-2",
        "name": "test-scenario-to-bulk-create-2-name",
        "enabled": true,
        "priority": 200001,
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
                "value": 3,
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
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "bulk-create-scenario-1",
        "status": 200,
        "item": {
          "_id": "bulk-create-scenario-1",
          "name": "test-scenario-to-bulk-create-1-name",
          "enabled": true,
          "priority": 200000,
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
                  "value": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ]
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "bulk-create-scenario-1",
          "name": "test-scenario-to-bulk-create-1-name",
          "enabled": true,
          "priority": 200000,
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
                  "value": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ]
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "errors": {
          "actions": "Actions is missing.",
          "enabled": "Enabled is missing.",
          "name": "Name is missing.",
          "triggers": "Triggers is missing."
        },
        "item": {}
      },
      {
        "status": 400,
        "errors": {
          "name": "Name already exists."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "priority": "Priority already exists."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.alarm_patterns": "AlarmPatterns is missing.",
          "actions.0.drop_scenario_if_not_matched": "DropScenarioIfNotMatched is missing.",
          "actions.0.emit_trigger": "EmitTrigger is missing.",
          "actions.0.entity_patterns": "EntityPatterns is missing.",
          "actions.0.type": "Type is missing."
        },
        "item": {
          "name": "test-scenario-to-bulk-create-4-name",
          "enabled": true,
          "priority": 13,
          "triggers": ["create"],
          "actions": [
            {
            }
          ]
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.alarm_patterns": "AlarmPatterns is missing.",
          "actions.0.entity_patterns": "EntityPatterns is missing."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.alarm_patterns": "Invalid alarm pattern list."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.alarm_patterns": "alarm pattern list contains an empty pattern.",
          "actions.0.entity_patterns": "entity pattern list contains an empty pattern."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.parameters.duration": "Duration is missing."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.parameters.ticket": "Ticket is missing."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.parameters.state": "State is missing."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.parameters.name": "Name is missing.",
          "actions.0.parameters.reason": "Reason doesn't exist.",
          "actions.0.parameters.type": "Type doesn't exist.",
          "actions.0.parameters.start_on_trigger": "StartOnTrigger or Tstart is required.",
          "actions.0.parameters.tstart": "Tstart or StartOnTrigger is required."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "actions.0.parameters.request": "Request is missing."
        },
        "item": {
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
        }
      },
      {
        "status": 400,
        "errors": {
          "_id": "ID already exists."
        },
        "item": {
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
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "id": "bulk-create-scenario-2",
        "status": 200,
        "item": {
          "_id": "bulk-create-scenario-2",
          "name": "test-scenario-to-bulk-create-2-name",
          "enabled": true,
          "priority": 200001,
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
                  "value": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ]
        }
      }
    ]
    """
    When I do GET /api/v4/scenarios?search=test-scenario-to-bulk-create&sort=asc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "bulk-create-scenario-1",
          "name": "test-scenario-to-bulk-create-1-name",
          "author": "root",
          "enabled": true,
          "priority": 200000,
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
                "duration": {
                  "value": 3,
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
          "_id": "bulk-create-scenario-2",
          "name": "test-scenario-to-bulk-create-2-name",
          "author": "root",
          "enabled": true,
          "priority": 200001,
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
                "duration": {
                  "value": 3,
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
