Feature: Bulk update a scenario
  I need to be able to bulk update scenarios
  Only admin should be able to bulk update scenarios

  Scenario: given bulk update request and no auth user should not allow access
    When I do PUT /api/v4/bulk/scenarios
    Then the response code should be 401

  Scenario: given bulk update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/scenarios
    Then the response code should be 403


  Scenario: given bulk update request should return multi status and should be handled independently
    When I am admin
    When I do PUT /api/v4/bulk/scenarios:
    """json
    [
      {
        "_id": "test-scenario-to-bulk-update-1",
        "name": "test-scenario-to-bulk-update-1-name",
        "enabled": true,
        "priority": 200010,
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
        "_id": "test-scenario-to-bulk-update-1",
        "name": "test-scenario-to-bulk-update-1-name-twice",
        "enabled": true,
        "priority": 200010,
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
        "priority": 12345
      },
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
      [],
      {
        "_id": "test-scenario-to-bulk-update-2",
        "name": "test-scenario-to-bulk-update-2-name",
        "enabled": true,
        "priority": 200011,
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
        "id": "test-scenario-to-bulk-update-1",
        "status": 200,
        "item": {
          "_id": "test-scenario-to-bulk-update-1",
          "name": "test-scenario-to-bulk-update-1-name",
          "enabled": true,
          "priority": 200010,
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
        "id": "test-scenario-to-bulk-update-1",
        "status": 200,
        "item": {
          "_id": "test-scenario-to-bulk-update-1",
          "name": "test-scenario-to-bulk-update-1-name-twice",
          "enabled": true,
          "priority": 200010,
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
        "errors": {
          "actions": "Actions is missing.",
          "enabled": "Enabled is missing.",
          "name": "Name is missing.",
          "triggers": "Triggers is missing.",
          "_id": "ID is missing."
        },
        "item": {
          "priority": 12345
        }
      },
      {
        "status": 400,
        "errors": {
          "name": "Name already exists.",
          "_id": "ID is missing."
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
          "priority": "Priority already exists.",
          "_id": "ID is missing."
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
          "actions.0.type": "Type is missing.",
          "_id": "ID is missing."
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
          "actions.0.entity_patterns": "EntityPatterns is missing.",
          "_id": "ID is missing."
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
          "actions.0.alarm_patterns": "Invalid alarm pattern list.",
          "_id": "ID is missing."
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
          "actions.0.entity_patterns": "entity pattern list contains an empty pattern.",
          "_id": "ID is missing."
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
          "actions.0.parameters.duration": "Duration is missing.",
          "_id": "ID is missing."
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
          "actions.0.parameters.ticket": "Ticket is missing.",
          "_id": "ID is missing."
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
          "actions.0.parameters.state": "State is missing.",
          "_id": "ID is missing."
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
          "actions.0.parameters.tstart": "Tstart or StartOnTrigger is required.",
          "_id": "ID is missing."
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
          "actions.0.parameters.request": "Request is missing.",
          "_id": "ID is missing."
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
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "id": "test-scenario-to-bulk-update-2",
        "status": 200,
        "item": {
          "_id": "test-scenario-to-bulk-update-2",
          "name": "test-scenario-to-bulk-update-2-name",
          "enabled": true,
          "priority": 200011,
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
    When I do GET /api/v4/scenarios?search=test-scenario-to-bulk-update&sort=asc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-scenario-to-bulk-update-1",
          "name": "test-scenario-to-bulk-update-1-name-twice",
          "author": "root",
          "enabled": true,
          "priority": 200010,
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
                "output": "test snooze updated",
                "duration": {
                  "value": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ],
          "created": 1605263992
        },
        {
          "_id": "test-scenario-to-bulk-update-2",
          "name": "test-scenario-to-bulk-update-2-name",
          "author": "root",
          "enabled": true,
          "priority": 200011,
          "delay": null,
          "disable_during_periods": null,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_patterns": [
                {
                  "_id": "test-scenario-to-bulk-update-2-alarm-updated"
                }
              ],
              "entity_patterns": null,
              "type": "snooze",
              "parameters": {
                "output": "test snooze updated",
                "duration": {
                  "value": 3,
                  "unit": "s"
                }
              },
              "drop_scenario_if_not_matched": false,
              "emit_trigger": false
            }
          ],
          "created": 1605263992
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
