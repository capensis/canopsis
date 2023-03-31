Feature: get running instruction
  I need to be able to get a instruction
  Only admin should be able to get a instruction

  Scenario: given started instruction should get it
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-get-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-get-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-get-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-get-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-get-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-get-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"m"},
              "description": "test-instruction-execution-get-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-get-1-step-2-operation-1",
              "time_to_complete": {"value": 2, "unit":"h"},
              "description": "test-instruction-execution-get-1-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-get-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-get-1-name",
      "description": "test-instruction-execution-get-1-description",
      "steps": [
        {
          "name": "test-instruction-execution-get-1-step-1",
          "time_to_complete": {"value": 181, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
              "name": "test-instruction-execution-get-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-1-step-1-operation-1-description"
            },
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-get-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"m"},
              "description": "test-instruction-execution-get-1-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-1-step-2",
          "time_to_complete": {"value": 2, "unit":"h"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-get-1-step-2-operation-1",
              "time_to_complete": {"value": 2, "unit":"h"},
              "description": "test-instruction-execution-get-1-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-1-step-2-endpoint"
        }
      ]
    }
    """
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"

  Scenario: given moved to next step instruction should get it
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-get-2-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-get-resource-2"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-get-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-get-2-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-get-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-get-2-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-2-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-get-2-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-2-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-get-2",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next
    Then the response code should be 200
    When I do PUT /api/v4/cat/executions/{{ .lastResponse._id }}/next-step
    Then the response code should be 200
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-get-2-name",
      "description": "test-instruction-execution-get-2-description",
      "steps": [
        {
          "name": "test-instruction-execution-get-2-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-get-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-1-operation-1-description",
              "jobs": []
            },
            {
              "name": "test-instruction-execution-get-2-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-2-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-2-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
              "name": "test-instruction-execution-get-2-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-get-2-step-2-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-2-step-2-endpoint"
        }
      ]
    }
    """
    When I save response step1CompletedAt={{ (index .lastResponse.steps 0).completed_at }}
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response op2CompletedAt={{ (index (index .lastResponse.steps 0).operations 1).completed_at }}
    When I save response op3StartedAt={{ (index (index .lastResponse.steps 1).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "op1StartedAt" >= "expectedStartedAt"
    Then "op1CompletedAt" >= "op1StartedAt"
    Then "op2StartedAt" >= "op1CompletedAt"
    Then "op2CompletedAt" >= "op2StartedAt"
    Then "step1CompletedAt" >= "op2CompletedAt"
    Then "op3StartedAt" >= "op2CompletedAt"

  Scenario: given moved to previous step instruction should get it
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-get-3-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-get-resource-3"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-get-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-get-3-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-get-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-3-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-get-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-get-3-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-3-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-get-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-get-1-step-3-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-get-3-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-get-3",
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
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-get-3-name",
      "description": "test-instruction-execution-get-3-description",
      "steps": [
        {
          "name": "test-instruction-execution-get-3-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "name": "test-instruction-execution-get-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-get-3-step-1-operation-1-description"
            },
            {
              "completed_at": null,
              "name": "test-instruction-execution-get-3-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-get-3-step-1-operation-2-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-3-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-get-3-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-get-3-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-get-1-step-3-operation-1-description",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-get-3-step-2-endpoint"
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

  Scenario: given unauth request should not allow access
    When I do GET /api/v4/cat/executions/notexist
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/executions/notexist
    Then the response code should be 403
