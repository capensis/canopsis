Feature: pause a instruction execution
  I need to be able to pause instruction operation
  Only admin should be able to pause a instruction

  Scenario: given running instruction should pause execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-resume-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-resume-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-resume-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-resume-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-resume-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-1-description"
            },
            {
              "name": "test-instruction-execution-resume-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-resume-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-resume-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-resume-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-2-operation-1-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-resume-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-resume-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    Then I save response executionID={{ .lastResponse._id }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/next
    Then the response code should be 200
    When I save response op1StartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response op1CompletedAt={{ (index (index .lastResponse.steps 0).operations 0).completed_at }}
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/pause
    Then the response code should be 204
    When I do PUT /api/v4/cat/executions/{{ .executionID }}/resume
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "description": "test-instruction-execution-resume-1-description",
      "name": "test-instruction-execution-resume-1-name",
      "status": 0,
      "steps": [
        {
          "endpoint": "test-instruction-execution-resume-1-step-1-endpoint",
          "name": "test-instruction-execution-resume-1-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "operations": [
            {
              "started_at": {{ .op1StartedAt }},
              "completed_at": {{ .op1CompletedAt }},
              "name": "test-instruction-execution-resume-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-1-description",
              "jobs": []
            },
            {
              "completed_at": null,
              "name": "test-instruction-execution-resume-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-1-operation-2-description",
              "jobs": []
            }
          ]
        },
        {
          "name": "test-instruction-execution-resume-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "endpoint": "test-instruction-execution-resume-1-step-2-endpoint",
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-execution-resume-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-resume-1-step-2-operation-1-description",
              "jobs": []
            }
          ]
        }
      ]
    }
    """
    When I save response op2StartedAt={{ (index (index .lastResponse.steps 0).operations 1).started_at }}
    When I save response expectedStartedAt=1
    Then "op2StartedAt" >= "expectedStartedAt"

  Scenario: given unauth request should not allow access
    When I do PUT /api/v4/cat/executions/notexist/resume
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/executions/notexist/resume
    Then the response code should be 403
