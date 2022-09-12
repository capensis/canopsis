Feature: move a instruction execution to next step
  I need to be able to run next instruction step
  Only admin should be able to run a instruction

  Scenario: given running instruction should complete current step and start next step of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-next-step-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-next-step-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-next-step-1-name",
      "description": "test-instruction-execution-next-step-1-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-1-step-1",
          "time_to_complete": {"value": 1, "unit":"s"},
          "failed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-1-operation-1-description connector test-instruction-execution-next-step-connector entity test-instruction-execution-next-step-resource-1/test-instruction-execution-next-step-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "completed_at": 0,
              "name": "test-instruction-execution-next-step-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-2-operation-1-description connector test-instruction-execution-next-step-connector entity test-instruction-execution-next-step-resource-1/test-instruction-execution-next-step-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.1.operations.0.started_at" should not be "0"

  Scenario: given running instruction should fail current step and start next step of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-2-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-next-step-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-2-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": false,
          "endpoint": "test-instruction-execution-next-step-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-2-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-2-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-next-step-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step:
    """json
    {
      "failed": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-next-step-2-name",
      "description": "test-instruction-execution-next-step-2-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-2-step-1",
          "time_to_complete": {"value": 1, "unit":"s"},
          "completed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-2-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "completed_at": 0,
              "name": "test-instruction-execution-next-step-2-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-2-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.failed_at" should not be "0"
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.1.operations.0.started_at" should not be "0"

  Scenario: given running instruction should complete execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-3-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-next-step-resource-3"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-3-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-3-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-next-step-3",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 5,
      "name": "test-instruction-execution-next-step-3-name",
      "description": "test-instruction-execution-next-step-3-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-3-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "failed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-3-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "failed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-3-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.1.started_at" should not be "0"
    Then the response key "steps.0.operations.1.completed_at" should not be "0"
    Then the response key "steps.1.completed_at" should not be "0"
    Then the response key "steps.1.operations.0.started_at" should not be "0"
    Then the response key "steps.1.operations.0.completed_at" should not be "0"

  Scenario: given running instruction should fail execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-5-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-next-step-resource-5"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-5-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-5-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-5-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-5-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-5-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-5-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-next-step-5",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step:
    """json
    {
      "failed": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 4,
      "name": "test-instruction-execution-next-step-5-name",
      "description": "test-instruction-execution-next-step-5-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-5-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-5-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-5-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-5-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "started_at": 0,
              "completed_at": 0,
              "name": "test-instruction-execution-next-step-5-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-5-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.failed_at" should not be "0"
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.1.started_at" should not be "0"
    Then the response key "steps.0.operations.1.completed_at" should not be "0"

  Scenario: given running instruction should fail execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-6-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-next-step-resource-6"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-6-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-6-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-6-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-6-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-6-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-6-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-6-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-next-step-6",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step:
    """json
    {
      "failed": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 4,
      "name": "test-instruction-execution-next-step-6-name",
      "description": "test-instruction-execution-next-step-6-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-6-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "failed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-6-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-6-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-6-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-6-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": 0,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-6-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-6-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-6-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.0.started_at" should not be "0"
    Then the response key "steps.0.operations.0.completed_at" should not be "0"
    Then the response key "steps.0.operations.1.started_at" should not be "0"
    Then the response key "steps.0.operations.1.completed_at" should not be "0"
    Then the response key "steps.1.failed_at" should not be "0"
    Then the response key "steps.1.operations.0.started_at" should not be "0"
    Then the response key "steps.1.operations.0.completed_at" should not be "0"

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/next-step
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/next-step
    Then the response code should be 403
