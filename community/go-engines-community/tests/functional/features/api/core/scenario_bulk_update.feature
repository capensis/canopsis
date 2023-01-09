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
        "priority": 7,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-1-alarm-updated"
                  }
                }
              ]
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
        "priority": 7,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-1-alarm-updated"
                  }
                }
              ]
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
      {},
      {
        "name": "test-scenario-to-check-unique-name-name",
        "enabled": true,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "triggers": ["create"],
        "actions": [{}]
      },
      {
        "name": "test-scenario-to-bulk-create-4-name",
        "enabled": true,
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "triggers": ["create"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-alarm"
                  }
                }
              ]
            ],
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-4-resource"
                  }
                }
              ]
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
        "priority": 8,
        "triggers": ["create","pbhenter"],
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-2-alarm-updated"
                  }
                }
              ]
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
        "_id": "test-scenario-to-bulk-update-3",
        "name": "test-scenario-to-bulk-update-3-name",
        "enabled": true,
        "priority": 17,
        "triggers": ["create"],
        "delay": {
          "value": 3,
          "unit": "s"
        },
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "v.component",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-3-alarm"
                  }
                }
              ]
            ],
            "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
          },
          {
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-3-name"
                  }
                }
              ]
            ],
            "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
          },
          {
            "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
            "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
      },
      {
        "name": "test-scenario-to-bulk-update-3-name",
        "enabled": true,
        "triggers": ["create"],
        "delay": {
          "value": 3,
          "unit": "s"
        },
        "actions": [
          {
            "alarm_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-3-alarm"
                  }
                }
              ]
            ],
            "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
          },
          {
            "entity_pattern": [
              [
                {
                  "field": "test",
                  "cond": {
                    "type": "eq",
                    "value": "test-scenario-to-bulk-update-3-name"
                  }
                }
              ]
            ],
            "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
          "priority": 7,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-1-alarm-updated"
                    }
                  }
                ]
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
          "priority": 7,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-1-alarm-updated"
                    }
                  }
                ]
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
        "item": {}
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "actions.0.drop_scenario_if_not_matched": "DropScenarioIfNotMatched is missing.",
          "actions.0.emit_trigger": "EmitTrigger is missing.",
          "actions.0.type": "Type is missing.",
          "_id": "ID is missing."
        },
        "item": {
          "name": "test-scenario-to-bulk-create-4-name",
          "enabled": true,
          "triggers": ["create"],
          "actions": [{}]
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "triggers": ["create"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-alarm"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-4-resource"
                    }
                  }
                ]
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
          "priority": 8,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-2-alarm-updated"
                    }
                  }
                ]
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
        "id": "test-scenario-to-bulk-update-3",
        "status": 200,
        "item": {
          "_id": "test-scenario-to-bulk-update-3",
          "name": "test-scenario-to-bulk-update-3-name",
          "enabled": true,
          "priority": 17,
          "triggers": ["create"],
          "delay": {
            "value": 3,
            "unit": "s"
          },
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-alarm"
                    }
                  }
                ]
              ],
              "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
            },
            {
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-name"
                    }
                  }
                ]
              ],
              "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
            },
            {
              "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
              "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
      },
      {
        "status": 400,
        "item": {
          "name": "test-scenario-to-bulk-update-3-name",
          "enabled": true,
          "triggers": ["create"],
          "delay": {
            "value": 3,
            "unit": "s"
          },
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-alarm"
                    }
                  }
                ]
              ],
              "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
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
            },
            {
              "entity_pattern": [
                [
                  {
                    "field": "test",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-name"
                    }
                  }
                ]
              ],
              "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
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
        },
        "errors": {
          "actions.0.alarm_pattern": "AlarmPattern is invalid alarm pattern.",
          "actions.1.entity_pattern": "EntityPattern is invalid entity pattern."
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
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "priority": 7,
          "delay": null,
          "disable_during_periods": null,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-1-alarm-updated"
                    }
                  }
                ]
              ],
              "old_alarm_patterns": null,
              "old_entity_patterns": null,
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
          "_id": "test-scenario-to-bulk-update-2",
          "name": "test-scenario-to-bulk-update-2-name",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "priority": 8,
          "delay": null,
          "disable_during_periods": null,
          "triggers": ["create","pbhenter"],
          "actions": [
            {
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-2-alarm-updated"
                    }
                  }
                ]
              ],
              "old_alarm_patterns": null,
              "old_entity_patterns": null,
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
          "_id": "test-scenario-to-bulk-update-3",
          "name": "test-scenario-to-bulk-update-3-name",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "enabled": true,
          "triggers": ["create"],
          "delay": {
            "value": 3,
            "unit": "s"
          },
          "actions": [
            {
              "old_alarm_patterns": null,
              "old_entity_patterns": null,
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-alarm"
                    }
                  }
                ]
              ],
              "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
              "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-pattern-to-rule-edit-2-pattern"
                    }
                  }
                ]
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
            },
            {
              "old_alarm_patterns": null,
              "old_entity_patterns": null,
              "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
              "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-pattern-to-rule-edit-1-pattern"
                    }
                  }
                ]
              ],
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-scenario-to-bulk-update-3-name"
                    }
                  }
                ]
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
            },
            {
              "old_alarm_patterns": null,
              "old_entity_patterns": null,
              "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
              "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-pattern-to-rule-edit-1-pattern"
                    }
                  }
                ]
              ],
              "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
              "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
              "entity_pattern": [
                [
                  {
                    "field": "name",
                    "cond": {
                      "type": "eq",
                      "value": "test-pattern-to-rule-edit-2-pattern"
                    }
                  }
                ]
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
