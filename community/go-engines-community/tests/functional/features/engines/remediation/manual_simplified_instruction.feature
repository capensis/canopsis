Feature: run an manual simplified instruction
  I need to be able to run an manual simplified instruction

  Scenario: given new alarm should run manual simplified instructions
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-1",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-1",
      "resource": "test-resource-to-run-manual-simplified-instruction-1",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-1
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-instruction-to-run-manual-simplified-instruction-1-name",
      "description": "test-instruction-to-run-manual-simplified-instruction-1-description",
      "status": 0,
      "instruction_type": 2,
      "complete_time": null,
      "completed_at": null,
      "jobs": [
        {
          "job_id": "test-job-to-run-manual-simplified-instruction-1",
          "name": "test-job-to-run-manual-simplified-instruction-1-name",
          "fail_reason": "",
          "status": 0,
          "started_at": null,
          "launched_at": null,
          "completed_at": null,
          "queue_number": 0
        },
        {
          "job_id": "test-job-to-run-manual-simplified-instruction-2",
          "name": "test-job-to-run-manual-simplified-instruction-2-name",
          "fail_reason": "",
          "status": 0,
          "started_at": null,
          "launched_at": null,
          "completed_at": null,
          "queue_number": null
        }
      ]
    }
    """
    When I save response startedAt={{ .lastResponse.started_at }}
    When I save response expectedStartedAt=1
    Then "startedAt" >= "expectedStartedAt"
    When I do GET /api/v4/cat/executions/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-instruction-to-run-manual-simplified-instruction-1-name",
      "description": "test-instruction-to-run-manual-simplified-instruction-1-description",
      "status": 0,
      "instruction_type": 2,
      "complete_time": null,
      "completed_at": null,
      "jobs": [
        {
          "job_id": "test-job-to-run-manual-simplified-instruction-1",
          "name": "test-job-to-run-manual-simplified-instruction-1-name"
        },
        {
          "job_id": "test-job-to-run-manual-simplified-instruction-2",
          "name": "test-job-to-run-manual-simplified-instruction-2-name"
        }
      ]
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-1-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction, instruction jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-2",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-2",
      "resource": "test-resource-to-run-manual-simplified-instruction-2",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-2
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-2"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-2-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction, instruction jobs shouldn't stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-3",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-3",
      "resource": "test-resource-to-run-manual-simplified-instruction-3",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-3
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-3"
    }
    """
    Then the response code should be 200
    When I wait the end of 3 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-3-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name. Job test-job-to-run-manual-simplified-instruction-2-name."
      },
      {
        "_t": "instructioncomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-3-name."
      }
    ]
    """

  Scenario: given new alarm should run manual instruction with failed jobs, instruction jobs should be stopped because of stopOnFail flag
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-4",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-4",
      "resource": "test-resource-to-run-manual-simplified-instruction-4",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-4
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-4"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "steps": {
          "page": 1,
          "limit": 20
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain:
    """json
    [
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionjobfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name. Job test-job-to-instruction-edit-1-name."
      },
      {
        "_t": "instructionfail",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-4-name."
      }
    ]
    """

  Scenario: given manual simplified instructions should not add steps to resolved alarm and new alarm
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-5",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-5",
      "resource": "test-resource-to-run-manual-simplified-instruction-5",
      "state": 1,
      "output": "test-output-to-run-manual-simplified-instruction-5"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-5
    Then the response code should be 200
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID }}",
      "instruction": "test-instruction-to-run-manual-simplified-instruction-5"
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-5",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-5",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-to-run-manual-simplified-instruction-5",
      "resource": "test-resource-to-run-manual-simplified-instruction-5"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-5",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-5",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-to-run-manual-simplified-instruction-5",
      "resource": "test-resource-to-run-manual-simplified-instruction-5"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-to-run-manual-simplified-instruction-5",
      "connector_name": "test-connector-name-to-run-manual-simplified-instruction-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-to-run-manual-simplified-instruction-5",
      "resource": "test-resource-to-run-manual-simplified-instruction-5",
      "state": 2,
      "output": "test-output-to-run-manual-simplified-instruction-5"
    }
    """
    When I wait the end of event processing
    When I wait 5s
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmID }}",
        "opened": false,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response array key "0.data.steps.data" should contain only:
    """json
    [
      {
        "_t": "stateinc",
        "val": 1
      },
      {
        "_t": "statusinc",
        "val": 1
      },
      {
        "_t": "instructionstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-5-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-5-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobcomplete",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-5-name. Job test-job-to-run-manual-simplified-instruction-1-name."
      },
      {
        "_t": "instructionjobstart",
        "a": "root",
        "user_id": "root",
        "m": "Instruction test-instruction-to-run-manual-simplified-instruction-5-name. Job test-job-to-run-manual-simplified-instruction-4-name."
      },
      {
        "_t": "cancel"
      },
      {
        "_t": "statusinc",
        "val": 4
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-to-run-manual-simplified-instruction-5&opened=true
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "opened": true,
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
