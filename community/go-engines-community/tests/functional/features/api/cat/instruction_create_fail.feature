Feature: Create a instruction
  I need to be able to create a instruction
  Only admin should be able to create a instruction

  @concurrent
  Scenario: given create request with custom id that already exist should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "_id": "test-instruction-to-check-unique-id"
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

  @concurrent
  Scenario: given create request with name that already exist should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "name": "test-instruction-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request to create manual instruction should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "steps": "Steps is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request to create auto instruction should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "enabled": "Enabled is missing.",
        "jobs": "Jobs is missing.",
        "name": "Name is missing.",
        "triggers": "Triggers is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "jobs": [],
      "triggers": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs": "Jobs is missing.",
        "triggers": "Triggers is missing."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": [
        {
          "type": "notexist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "triggers.0.type": "Type must be one of [create statedec stateinc changestate unsnooze activate pbhenter pbhleave eventscount]."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request to create simplified manual instruction should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "enabled": "Enabled is missing.",
        "jobs": "Jobs is missing.",
        "name": "Name is missing.",
        "timeout_after_execution": "TimeoutAfterExecution is missing."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "jobs": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs": "Jobs is missing."
      }
    }
    """

  @concurrent
  Scenario: given create request with not exist job should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "steps": [
        {
          "operations": [
            {
              "jobs": [
                "test-job-not-exist"
              ]
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "steps.0.operations.0.jobs": "Jobs doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "jobs": [
        {
          "job": "test-job-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs.0.job": "Job doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 2,
      "jobs": [
        {
          "job": "test-job-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "jobs.0.job": "Job doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with not exist pbehavior type should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "active_on_pbh": ["notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "active_on_pbh": "ActiveOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "disabled_on_pbh": ["notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "disabled_on_pbh": "DisabledOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "active_on_pbh": ["test-default-maintenance-type", "notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "active_on_pbh": "ActiveOnPbh doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "disabled_on_pbh": ["test-default-maintenance-type", "notexist"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "disabled_on_pbh": "DisabledOnPbh doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with invalid patterns should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "entity_pattern": [
        []
      ],
      "alarm_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "v.created_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-instruction-to-create-fail-8-pattern"
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "corporate_alarm_pattern": "test-pattern-not-exist",
      "type": 1,
      "name": "test-instruction-to-create-fail-8-name",
      "description": "test-instruction-to-create-fail-8-description",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "job": "test-job-to-instruction-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_alarm_pattern": "CorporateAlarmPattern doesn't exist."
      }
    }
    """
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "corporate_entity_pattern": "test-pattern-not-exist",
      "type": 1,
      "name": "test-instruction-to-create-fail-8-name",
      "description": "test-instruction-to-create-fail-8-description",
      "enabled": true,
      "triggers": [
        {
          "type": "create"
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "job": "test-job-to-instruction-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create auto instruction requests with invalid threshold value should return errors
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-fail-9-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-fail-9-description",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount"
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": true,
          "job": "test-job-to-instruction-edit-1"
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
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-fail-9-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-fail-9-description",
      "enabled": true,
      "triggers": [
        {
          "type": "eventscount",
          "threshold": 1
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": true,
          "job": "test-job-to-instruction-edit-1"
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
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-to-create-fail-9-name",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "description": "test-instruction-to-create-fail-9-description",
      "enabled": true,
      "triggers": [
        {
          "type": "create",
          "threshold": 1
        }
      ],
      "timeout_after_execution": {
        "value": 10,
        "unit": "m"
      },
      "jobs": [
        {
          "stop_on_fail": true,
          "job": "test-job-to-instruction-edit-1"
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

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/instructions
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/instructions
    Then the response code should be 403
