Feature: Create a scenario
  I need to be able to create a scenario
  Only admin should be able to create a scenario

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/scenarios
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/scenarios
    Then the response code should be 403

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "_id": "test-scenario-to-create-1",
      "name": "test-scenario-to-create-1-name",
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
                  "value": "test-scenario-to-create-1-action-1-alarm"
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
                  "value": "test-scenario-to-create-1-action-1-resource"
                }
              }
            ]
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
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-1-action-2-alarm"
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
                  "value": "test-scenario-to-create-1-action-2-resource"
                }
              }
            ]
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
              "payload": "{\"test-scenario-to-create-1-action-2-payload\": \"test-scenario-to-create-1-action-2-paload-value\"}",
              "timeout": {
                "value": 1,
                "unit": "m"
              },
              "retry_count": 3,
              "retry_delay": {
                "value": 3,
                "unit": "s"
              }
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-1-action-2-ticket",
              "test-scenario-to-create-1-action-2-info": "test-scenario-to-create-1-action-2-info-value"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        },
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-1-action-3-alarm"
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
                  "value": "test-scenario-to-create-1-action-3-resource"
                }
              }
            ]
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
      "_id": "test-scenario-to-create-1",
      "name": "test-scenario-to-create-1-name",
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
                  "value": "test-scenario-to-create-1-action-1-alarm"
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
                  "value": "test-scenario-to-create-1-action-1-resource"
                }
              }
            ]
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
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-1-action-2-alarm"
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
                  "value": "test-scenario-to-create-1-action-2-resource"
                }
              }
            ]
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
              "payload": "{\"test-scenario-to-create-1-action-2-payload\": \"test-scenario-to-create-1-action-2-paload-value\"}",
              "timeout": {
                "value": 1,
                "unit": "m"
              },
              "retry_count": 3,
              "retry_delay": {
                "value": 3,
                "unit": "s"
              }
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-1-action-2-ticket",
              "test-scenario-to-create-1-action-2-info": "test-scenario-to-create-1-action-2-info-value"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
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
                  "value": "test-scenario-to-create-1-action-3-alarm"
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
                  "value": "test-scenario-to-create-1-action-3-resource"
                }
              }
            ]
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
              "priority": 25,
              "type": "maintenance"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    When I do GET /api/v4/scenarios/test-scenario-to-create-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-scenario-to-create-1",
      "name": "test-scenario-to-create-1-name",
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
                  "value": "test-scenario-to-create-1-action-1-alarm"
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
                  "value": "test-scenario-to-create-1-action-1-resource"
                }
              }
            ]
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
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-1-action-2-alarm"
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
                  "value": "test-scenario-to-create-1-action-2-resource"
                }
              }
            ]
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
              "payload": "{\"test-scenario-to-create-1-action-2-payload\": \"test-scenario-to-create-1-action-2-paload-value\"}",
              "timeout": {
                "value": 1,
                "unit": "m"
              },
              "retry_count": 3,
              "retry_delay": {
                "value": 3,
                "unit": "s"
              }
            },
            "declare_ticket": {
              "empty_response": false,
              "is_regexp": false,
              "ticket_id": "test-scenario-to-create-1-action-2-ticket",
              "test-scenario-to-create-1-action-2-info": "test-scenario-to-create-1-action-2-info-value"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
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
                  "value": "test-scenario-to-create-1-action-3-alarm"
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
                  "value": "test-scenario-to-create-1-action-3-resource"
                }
              }
            ]
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
              "priority": 25,
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
      "name": "test-scenario-to-create-2-name",
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
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-2-alarm"
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
                  "value": "test-scenario-to-create-2-resource"
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-scenario-to-create-2-name",
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
                  "value": "test-scenario-to-create-2-alarm"
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
                  "value": "test-scenario-to-create-2-resource"
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
    
  Scenario: given create request with corporate patterns in different actions should return success
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-3-name",
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
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-3-alarm"
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
                  "value": "test-scenario-to-create-3-name"
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
        },
        {
          "old_entity_patterns": [
            {
              "name": "test-scenario-to-update-3-name-1"
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
        },
        {
          "old_alarm_patterns": [
            {
              "v": {
                "component": "test-scenario-to-update-3-name-1"
              }
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
                  "value": "test-scenario-to-create-3-alarm"
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
                  "value": "test-scenario-to-create-3-name"
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
        },
        {
          "old_entity_patterns": [
            {
              "name": "test-scenario-to-update-3-name-1"
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
        },
        {
          "old_alarm_patterns": [
            {
              "v": {
                "component": "test-scenario-to-update-3-name-1"
              }
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

  Scenario: given create request with invalid patterns in different actions should return errors
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-to-create-4-name",
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
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/scenarios:
    """json
    {}
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-resource"
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
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-resource"
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
      "triggers": ["create"],
      "actions": [
        {
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "actions.0.drop_scenario_if_not_matched": "DropScenarioIfNotMatched is missing.",
        "actions.0.emit_trigger": "EmitTrigger is missing.",
        "actions.0.type": "Type is missing."
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
                }
              }
            ]
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
      "triggers": ["create"],
      "actions": [
        {
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-scenario-to-create-4-alarm"
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
                  "value": "test-scenario-to-create-4-name"
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
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """
