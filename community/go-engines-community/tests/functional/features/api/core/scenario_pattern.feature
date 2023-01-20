Feature: Update and delete corporate pattern should affect eventfilter models
  Scenario: given eventfilter and corporate pattern update and delete actions should update patterns in eventfiler
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-2",
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-1",
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/scenarios/scenario-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": ["create"],
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-2",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-2-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-2-pattern"
                }
              }
            ],
            [
              {
                "field": "v.creation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              },
              {
                "field": "v.activation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              },
              {
                "field": "v.ack.t",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-1",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-1-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-1-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-3-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-3-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-4-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-4-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-scenario-corporate-update-1:
    """json
    {
      "title": "new entity pattern title",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "new entity pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/scenarios/scenario-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": ["create"],
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-2",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-2-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-2-pattern"
                }
              }
            ],
            [
              {
                "field": "v.creation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              },
              {
                "field": "v.activation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              },
              {
                "field": "v.ack.t",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263992,
                    "to": 1605264992
                  }
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-1",
          "corporate_entity_pattern_title": "new entity pattern title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "new entity pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-3-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-3-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-4-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-4-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    When I do PUT /api/v4/patterns/test-pattern-to-scenario-corporate-update-2:
    """json
    {
      "title": "new alarm pattern title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "new alarm pattern"
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/scenarios/scenario-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": ["create"],
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-2",
          "corporate_alarm_pattern_title": "new alarm pattern title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "new alarm pattern"
                }
              }
            ],
            [
              {
                "field": "v.creation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.activation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.ack.t",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-1",
          "corporate_entity_pattern_title": "new entity pattern title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "new entity pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-3-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-3-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-4-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-4-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-scenario-corporate-update-1
    Then the response code should be 204
    When I do GET /api/v4/scenarios/scenario-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": ["create"],
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-2",
          "corporate_alarm_pattern_title": "new alarm pattern title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "new alarm pattern"
                }
              }
            ],
            [
              {
                "field": "v.creation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.activation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.ack.t",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "new entity pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-3-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-3-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-4-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-4-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    Then the response key "actions.3.corporate_entity_pattern" should not exist
    Then the response key "actions.3.corporate_entity_pattern_title" should not exist
    When I do DELETE /api/v4/patterns/test-pattern-to-scenario-corporate-update-2
    Then the response code should be 204
    When I do GET /api/v4/scenarios/scenario-pattern-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "scenario-pattern-1",
      "name": "scenario-pattern-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": ["create"],
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
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 1"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "simple pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 2"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "new alarm pattern"
                }
              }
            ],
            [
              {
                "field": "v.creation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.activation_date",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              },
              {
                "field": "v.ack.t",
                "cond": {
                  "type": "absolute_time",
                  "value": {
                    "from": 1605263993,
                    "to": 1605264993
                  }
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 3"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "new entity pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 4"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_entity_pattern": "test-pattern-to-scenario-corporate-update-3",
          "corporate_entity_pattern_title": "test-pattern-to-scenario-corporate-update-3-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-3-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 5"
        },
        {
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "corporate_alarm_pattern": "test-pattern-to-scenario-corporate-update-4",
          "corporate_alarm_pattern_title": "test-pattern-to-scenario-corporate-update-4-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-scenario-corporate-update-4-pattern"
                }
              }
            ]
          ],
          "type": "snooze",
          "parameters": {
            "author": "scenario-pattern-1-action-1-author",
            "output": "scenario-pattern-1-action-1-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "snooze 6"
        }
      ]
    }
    """
    Then the response key "actions.2.corporate_alarm_pattern" should not exist
    Then the response key "actions.2.corporate_alarm_pattern_title" should not exist
    Then the response key "actions.3.corporate_entity_pattern" should not exist
    Then the response key "actions.3.corporate_entity_pattern_title" should not exist
