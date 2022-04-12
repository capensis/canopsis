Feature: run a instruction
  I need to be able to run a instruction
  Only admin should be able to run a instruction

  Scenario: given instruction should start first operation of instruction
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-execution-start-1-name",
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-instruction-execution-start-resource-1"
            }
          }
        ]
      ],
      "description": "test-instruction-execution-start-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 10,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-execution-start-1-step-1",
          "operations": [
            {
              "name": "test-instruction-execution-start-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-start-1-step-1-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}",
              "jobs": ["test-instruction-execution-1"]
            },
            {
              "name": "test-instruction-execution-start-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "test-instruction-execution-start-1-step-1-operation-2-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-start-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-start-1-step-2",
          "operations": [
            {
              "name": "test-instruction-execution-start-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "test-instruction-execution-start-1-step-2-operation-1-description connector {{ `{{ .Alarm.Value.Connector }}` }} entity {{ `{{ .Entity.ID }}` }}"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-execution-start-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-1",
      "instruction": "{{ .lastResponse._id }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "status": 0,
      "name": "test-instruction-execution-start-1-name",
      "description": "test-instruction-execution-start-1-description",
      "steps": [
        {
          "name": "test-instruction-execution-start-1-step-1",
          "time_to_complete": {"value": 4, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "completed_at": 0,
              "name": "test-instruction-execution-start-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-execution-start-1-step-1-operation-1-description connector test-instruction-execution-start-connector entity test-instruction-execution-start-resource-1/test-instruction-execution-start-component",
              "jobs": [
                {
                  "_id": "",
                  "job_id": "test-instruction-execution-1",
                  "status": null,
                  "name": "test-instruction-execution-1-name",
                  "fail_reason": "",
                  "started_at": 0,
                  "launched_at": 0,
                  "completed_at": 0
                }
              ]
            },
            {
              "started_at": 0,
              "completed_at": 0,
              "name": "test-instruction-execution-start-1-step-1-operation-2",
              "time_to_complete": {"value": 3, "unit":"s"},
              "description": "",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-start-1-step-1-endpoint"
        },
        {
          "name": "test-instruction-execution-start-1-step-2",
          "time_to_complete": {"value": 6, "unit":"s"},
          "completed_at": 0,
          "failed_at": 0,
          "operations": [
            {
              "started_at": 0,
              "completed_at": 0,
              "name": "test-instruction-execution-start-1-step-2-operation-1",
              "time_to_complete": {"value": 6, "unit":"s"},
              "description": "",
              "jobs": []
            }
          ],
          "endpoint": "test-instruction-execution-start-1-step-2-endpoint"
        }
      ]
    }
    """
    Then the response key "steps.0.operations.0.started_at" should not be "0"

  Scenario: given instruction should not start instruction multiple times for one alarm
    When I am admin
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-2",
      "instruction": "test-instruction-execution-to-start-multiple-times-1"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-2",
      "instruction": "test-instruction-execution-to-start-multiple-times-1"
    }
    """
    Then the response code should be 400

  Scenario: given instruction should start instruction multiple times for multiple alarms
    When I am admin
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-3-1",
      "instruction": "test-instruction-execution-to-start-multiple-times-2"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-3-2",
      "instruction": "test-instruction-execution-to-start-multiple-times-2"
    }
    """
    Then the response code should be 200

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
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-1",
      "instruction": "test-instruction-not-exist"
    }
    """
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "test-instruction-execution-start-1",
      "instruction": "test-instruction-not-exist"
    }
    """
    Then the response code should be 403
