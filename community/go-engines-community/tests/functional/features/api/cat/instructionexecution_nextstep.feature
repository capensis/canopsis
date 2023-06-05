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
              "value": "test-resource-to-instruction-execution-next-step-1"
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
      "alarm": "test-alarm-instruction-execution-next-step-1",
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
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-1-operation-1-description connector test-connector-default entity test-resource-to-instruction-execution-next-step-1/test-component-default",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
              "name": "test-instruction-execution-next-step-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-1-step-2-operation-1-description connector test-connector-default entity test-resource-to-instruction-execution-next-step-1/test-component-default",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-1-step-2-endpoint"
        }
      ]
    }
    """
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "step1CompletedAt" >= "op1CompletedAt"
    Then "op2StartedAt" >= "op1CompletedAt"

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
              "value": "test-resource-to-instruction-execution-next-step-2"
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
      "alarm": "test-alarm-instruction-execution-next-step-2",
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
          "completed_at": null,
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
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
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
    When I save response step1FailedAt={{ (index .lastResponse.steps 0).failed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "step1FailedAt" >= "op1CompletedAt"
    Then "op2StartedAt" >= "op1CompletedAt"

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
              "value": "test-resource-to-instruction-execution-next-step-3"
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
      "alarm": "test-alarm-instruction-execution-next-step-3",
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
          "failed_at": null,
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
          "failed_at": null,
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
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response step2CompletedAt={{ (index .lastResponse.steps 1).completed_at }}
    When I save response op3StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response op3CompletedAt={{ (index (index .lastResponse.steps 1).operations 0).completed_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1CompletedAt" >= "op2CompletedAt"
    Then "op3StartedAt" >= "op2CompletedAt"
    Then "op3CompletedAt" >= "op3StartedAt"
    Then "step2CompletedAt" >= "op3CompletedAt"

  Scenario: given removed instruction should complete execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-4-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-instruction-execution-next-step-4"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-4-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-4-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-4-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-4-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-4-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-4-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-instruction-execution-next-step-4",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I do DELETE /api/v4/cat/instructions/{{ .instructionID }}
    Then the response code should be 204
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 5,
      "name": "test-instruction-execution-next-step-4-name",
      "description": "test-instruction-execution-next-step-4-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-4-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-4-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-4-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-4-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-4-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-4-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-4-step-2-endpoint"
        }
      ]
    }
    """
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response step2CompletedAt={{ (index .lastResponse.steps 1).completed_at }}
    When I save response op3StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response op3CompletedAt={{ (index (index .lastResponse.steps 1).operations 0).completed_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1CompletedAt" >= "op2CompletedAt"
    Then "op3StartedAt" >= "op2CompletedAt"
    Then "op3CompletedAt" >= "op3StartedAt"
    Then "step2CompletedAt" >= "op3CompletedAt"

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
              "value": "test-resource-to-instruction-execution-next-step-5"
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
      "alarm": "test-alarm-instruction-execution-next-step-5",
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
          "completed_at": null,
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
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-next-step-5-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-5-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-5-step-2-endpoint"
        }
      ]
    }
    """
    When I save response step1FailedAt={{ (index .lastResponse.steps 0).failed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1FailedAt" >= "op2CompletedAt"

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
              "value": "test-resource-to-instruction-execution-next-step-6"
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
      "alarm": "test-alarm-instruction-execution-next-step-6",
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
          "failed_at": null,
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
          "completed_at": null,
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
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response step2FailedAt={{ (index .lastResponse.steps 1).failed_at }}
    When I save response op3StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response op3CompletedAt={{ (index (index .lastResponse.steps 1).operations 0).completed_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1CompletedAt" >= "op2CompletedAt"
    Then "op3StartedAt" >= "op2CompletedAt"
    Then "op3CompletedAt" >= "op3StartedAt"
    Then "step2FailedAt" >= "op3CompletedAt"

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/next-step
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/next-step
    Then the response code should be 403

  Scenario: given updated instruction should complete execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-7-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-instruction-execution-next-step-7"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-7-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-7-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-7-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-7-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-7-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-7-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-instruction-execution-next-step-7",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I save response executionID={{ .lastResponse._id }}
    When I do PUT /api/v4/cat/instructions/{{ .instructionID }}:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-next-step-7-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-instruction-execution-next-step-7"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-next-step-7-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-next-step-7-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-next-step-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-1-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-next-step-7-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next-step
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 5,
      "name": "test-instruction-execution-next-step-7-name",
      "description": "test-instruction-execution-next-step-7-description",
      "steps": [
        {
          "name": "test-instruction-execution-next-step-7-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-7-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-next-step-7-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-7-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-next-step-7-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-next-step-7-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-next-step-7-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-next-step-7-step-2-endpoint"
        }
      ]
    }
    """
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response step2CompletedAt={{ (index .lastResponse.steps 1).completed_at }}
    When I save response op3StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response op3CompletedAt={{ (index (index .lastResponse.steps 1).operations 0).completed_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1CompletedAt" >= "op2CompletedAt"
    Then "op3StartedAt" >= "op2CompletedAt"
    Then "op3CompletedAt" >= "op3StartedAt"
    Then "step2CompletedAt" >= "op3CompletedAt"
