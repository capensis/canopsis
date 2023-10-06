Feature: Update a scenario
  I need to be able to update a scenario
  Only admin should be able to update a scenario

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/scenarios/notexist
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/scenarios/notexist
    Then the response code should be 403

  @concurrent
  Scenario: given update scenario request should return ok
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """json
    {
      "name": "test-scenario-to-update-1-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        },
        {
          "type": "pbhenter"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-1-alarm-updated"
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
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-to-update-1",
      "name": "test-scenario-to-update-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "delay": null,
      "disable_during_periods": null,
      "triggers": [
        {
          "type": "create"
        },
        {
          "type": "pbhenter"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-1-alarm-updated"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_alarm_patterns": null,
          "type": "snooze",
          "parameters": {
            "output": "test snooze updated",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """

  @concurrent
  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-1:
    """json
    {
      "name": "test-scenario-to-check-unique-name-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        },
        {
          "type": "pbhenter"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-1-alarm-updated"
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

  @concurrent
  Scenario: given no exist scenario id should return error
    When I am admin
    When I do PUT /api/v4/scenarios/notexist:
    """json
    {
      "name": "test-scenario-to-update-notexist-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-notexist-alarm"
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
    """
    Then the response code should be 404

  @concurrent
  Scenario: given create request with custom_id shouldn't update id
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-check-id:
    """json
    {
      "_id": "change-id",
      "name": "my_scenario-name-new",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-check-id"
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
                  "value": "test-scenario-to-check-id"
                }
              }
            ]
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-to-check-id",
      "name": "my_scenario-name-new"
    }
    """

  @concurrent
  Scenario: given update request with corporate patterns in different actions should return success
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-2:
    """json
    {
      "name": "test-scenario-to-update-2-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
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
                  "value": "test-scenario-to-update-2-alarm"
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
                  "value": "test-scenario-to-update-2-name"
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
    """   
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-update-2-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
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
                  "value": "test-scenario-to-update-2-alarm"
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
                  "value": "test-scenario-to-update-2-name"
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
    """

  @concurrent
  Scenario: given update request with invalid patterns in different actions should return errors
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-2:
    """json
    {
      "name": "test-scenario-to-update-2-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
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
                  "value": "test-scenario-to-update-2-alarm"
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
                  "value": "test-scenario-to-update-2-name"
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "actions.0.alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "actions.1.entity_pattern": "EntityPattern is invalid entity pattern.",
        "actions.2.alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "actions.2.entity_pattern": "EntityPattern or AlarmPattern is required."
      }
    }
    """

  @concurrent
  Scenario: given update request should update old patterns but keep them if they're not updated
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-3:
    """json
    {
      "name": "test-scenario-to-update-3-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
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
                  "value": "test-scenario-to-update-3-alarm"
                }
              }
            ]
          ],
          "old_entity_patterns": [
            {
              "name": "test-scenario-to-update-3-name-1"
            }
          ],
          "comment": "first ack",
          "type": "ack",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "old_alarm_patterns": [
            {
              "_id": "test-scenario-to-update-3-alarm-2"
            }
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-3-name"
                }
              }
            ]
          ],
          "comment": "second ack",
          "type": "ack",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-update-3-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
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
                  "value": "test-scenario-to-update-3-alarm"
                }
              }
            ]
          ],
          "old_entity_patterns": [
            {
              "name": "test-scenario-to-update-3-name-1"
            }
          ],
          "type": "ack",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "first ack"
        },
        {
          "old_alarm_patterns": [
            {
              "_id": "test-scenario-to-update-3-alarm-2"
            }
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-3-name"
                }
              }
            ]
          ],
          "type": "ack",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "second ack"
        }
      ]
    }    
    """    

  @concurrent
  Scenario: given update scenario request should return ok
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-4:
    """json
    {
      "name": "test-scenario-to-update-4-name",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 3
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-4-alarm-updated"
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
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-to-update-4",
      "name": "test-scenario-to-update-4-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": true,
      "delay": null,
      "disable_during_periods": null,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 3
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-4-alarm-updated"
                }
              }
            ]
          ],
          "old_entity_patterns": null,
          "old_alarm_patterns": null,
          "type": "snooze",
          "parameters": {
            "output": "test snooze updated",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    
  @concurrent
  Scenario: given update scenario request with invalid threshold value should return errors
    When I am admin
    When I do PUT /api/v4/scenarios/test-scenario-to-update-5:
    """json
    {
      "name": "test-scenario-to-update-5-name",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount"
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-5-alarm-updated"
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
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "triggers.0.threshold": "Threshold is required when Type eventscount is defined."
      }
    }
    """
    When I do PUT /api/v4/scenarios/test-scenario-to-update-5:
    """json
    {
      "name": "test-scenario-to-update-5-name",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 1
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-5-alarm-updated"
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
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "triggers.0.threshold": "Threshold should be greater than 1."
      }
    }
    """
    When I do PUT /api/v4/scenarios/test-scenario-to-update-5:
    """json
    {
      "name": "test-scenario-to-update-5-name",
      "enabled": true,
      "triggers": [
        {
          "type": "create",
          "threshold": 3
        }
      ],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-update-5-alarm-updated"
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
          "emit_trigger": false,
          "comment": "test comment"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "triggers.0.threshold": "Threshold should be empty when Type eventscount is not defined."
      }
    }
    """
