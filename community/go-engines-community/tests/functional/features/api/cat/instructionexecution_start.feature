Feature: run a instruction
  I need to be able to run a instruction
  Only admin should be able to run a instruction

  Scenario: given instruction should start first operation of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-execution-start-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-instruction-execution-start-1"
            }
          }
        ]
      ],
      "description": "test-instruction-to-execution-start-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-execution-start-1-step-1",
          "operations": [
            {
              "name": "test-instruction-to-execution-start-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}",
              "jobs": ["test-instruction-execution-1"]
            },
            {
              "name": "test-instruction-to-execution-start-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-1-operation-2-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-execution-start-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-execution-start-1-step-2",
          "operations": [
            {
              "name": "test-instruction-to-execution-start-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-execution-start-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-to-execution-start-1-name",
      "description": "test-instruction-to-execution-start-1-description",
      "instruction_type": 0,
      "completed_at": null,
      "complete_time": null,
      "steps": [
        {
          "name": "test-instruction-to-execution-start-1-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "completed_at": null,
              "name": "test-instruction-to-execution-start-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-1-operation-1-description connector test-connector-default entity test-resource-to-instruction-execution-start-1/test-component-default",
              "jobs": [
                {
                  "_id": "",
                  "job_id": "test-instruction-execution-1",
                  "status": null,
                  "name": "test-instruction-execution-1-name",
                  "fail_reason": "",
                  "queue_number": null,
                  "started_at": null,
                  "launched_at": null,
                  "completed_at": null
                }
              ]
            },
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-to-execution-start-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-1-operation-2-description connector test-connector-default entity test-resource-to-instruction-execution-start-1/test-component-default",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-to-execution-start-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-to-execution-start-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": null,
          "failed_at": null,
          "operations": [
            {
              "started_at": null,
              "completed_at": null,
              "name": "test-instruction-to-execution-start-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-to-execution-start-1-step-2-operation-1-description connector test-connector-default entity test-resource-to-instruction-execution-start-1/test-component-default",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-to-execution-start-1-step-2-endpoint"
        }
      ]
    }
    """
    When I save response startedAt={{ .lastResponse.started_at }}
    When I save response operationStartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "startedAt" >= "expectedStartedAt"
    Then "operationStartedAt" >= "expectedStartedAt"

  Scenario: given instruction should not start instruction multiple times for one alarm
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-execution-start-2-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-instruction-execution-start-2"
            }
          }
        ]
      ],
      "description": "test-instruction-to-execution-start-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-execution-start-2-step-1",
          "operations": [
            {
              "name": "test-instruction-to-execution-start-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-execution-start-2-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}",
              "jobs": ["test-instruction-execution-1"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-execution-start-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-2",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-2",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 400

  Scenario: given instruction should start instruction multiple times for multiple alarms
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-execution-start-3-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-instruction-execution-start-3-1",
                "test-resource-to-instruction-execution-start-3-2"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-to-execution-start-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-execution-start-3-step-1",
          "operations": [
            {
              "name": "test-instruction-to-execution-start-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-execution-start-3-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}",
              "jobs": ["test-instruction-execution-1"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-execution-start-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-3-1",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-3-2",
      "instruction": "{{ .instructionID }}"
    }
    """
    Then the response code should be 200

  Scenario: given instruction with old pattern should start an execution
    When I am admin
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-4-1",
      "instruction": "test-instruction-to-execution-start-4-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-instruction-to-execution-start-4-1-name"
    }
    """
    When I save response startedAt={{ .lastResponse.started_at }}
    When I save response operationStartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "startedAt" >= "expectedStartedAt"
    Then "operationStartedAt" >= "expectedStartedAt"
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-4-2",
      "instruction": "test-instruction-to-execution-start-4-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-instruction-to-execution-start-4-2-name"
    }
    """
    When I save response startedAt={{ .lastResponse.started_at }}
    When I save response operationStartedAt={{ (index (index .lastResponse.steps 0).operations 0).started_at }}
    When I save response expectedStartedAt=1
    Then "startedAt" >= "expectedStartedAt"
    Then "operationStartedAt" >= "expectedStartedAt"

  Scenario: given instruction with empty patterns should not start an execution
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-to-execution-start-5-name",
      "description": "test-instruction-to-execution-start-5-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-to-execution-start-5-step-1",
          "operations": [
            {
              "name": "test-instruction-to-execution-start-5-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-to-execution-start-5-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}",
              "jobs": ["test-instruction-execution-1"]
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-to-execution-start-5-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-instruction-execution-start-5",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 400

  Scenario: given invalid request should return errors
    When I am admin
    When I do POST /api/v4/cat/executions
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm": "Alarm is missing.",
        "instruction": "Instruction is missing."
      }
    }
    """

  Scenario: given auto instruction should return errors
    When I am admin
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-alarm-to-not-run-auto-instruction-manually",
      "instruction": "test-instruction-to-not-run-auto-instruction-manually"
    }
    """
    Then the response code should be 404

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/cat/executions
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/executions
    Then the response code should be 403
