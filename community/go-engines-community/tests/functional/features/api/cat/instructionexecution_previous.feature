Feature: move a instruction execution to previous operation
  I need to be able to run previous instruction operation
  Only admin should be able to run a instruction

  Scenario: given running instruction should start previous operation of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-previous-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-previous-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-previous-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-previous-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-previous-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            },
            {
              "name": "test-instruction-execution-previous-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-1-operation-2-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-previous-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-previous-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-previous-1-name",
      "description": "test-instruction-execution-previous-1-description",
      "steps": [
        {
          "name": "test-instruction-execution-previous-1-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
              "name": "test-instruction-execution-previous-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-1-operation-1-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-1/test-instruction-execution-previous-component"
            },
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-previous-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-1-operation-2-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-1/test-instruction-execution-previous-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-previous-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-1-step-2-operation-1-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-1/test-instruction-execution-previous-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-1-step-2-endpoint"
        }
      ]
    }
    """
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"

  Scenario: given running instruction with first operation started should return error
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-previous-2-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-previous-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-previous-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-previous-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-previous-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-previous-2-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-2-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-previous-2-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-previous-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
    Then the response code should be 400

  Scenario: given running instruction should start previous operation of previous step
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-previous-3-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-previous-resource-3"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-previous-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-previous-3-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-previous-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            },
            {
              "name": "test-instruction-execution-previous-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-1-operation-2-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-3-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-previous-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-previous-3",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-previous-3-name",
      "description": "test-instruction-execution-previous-3-description",
      "steps": [
        {
          "name": "test-instruction-execution-previous-3-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-previous-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-1-operation-1-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-3/test-instruction-execution-previous-component"
            },
            {
              "completed_at": null,
              "name": "test-instruction-execution-previous-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-1-operation-2-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-3/test-instruction-execution-previous-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-3-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-previous-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-3-step-2-operation-1-description connector test-instruction-execution-previous-connector entity test-instruction-execution-previous-resource-3/test-instruction-execution-previous-component",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-3-step-2-endpoint"
        }
      ]
    }
    """
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"

  Scenario: given running instruction should finish execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-previous-4-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-previous-resource-4"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-previous-4-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-previous-4-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-previous-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-previous-4-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-4-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-4-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-previous-4-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-previous-4-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-previous-4",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/previous
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
      "name": "test-instruction-execution-previous-4-name",
      "description": "test-instruction-execution-previous-4-description",
      "steps": [
        {
          "name": "test-instruction-execution-previous-4-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "operations": [
            {
              "name": "test-instruction-execution-previous-4-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-previous-4-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-4-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-previous-4-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "operations": [
            {
              "name": "test-instruction-execution-previous-4-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-previous-4-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-previous-4-step-2-endpoint"
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

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/previous
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/previous
    Then the response code should be 403
