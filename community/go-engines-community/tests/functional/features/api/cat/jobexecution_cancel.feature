Feature: cancel a job execution
  I need to be able to cancel a instruction execution
  Only admin should be able to cancel a job execution

  Scenario: given job should cancel running executions on job update
    When I am admin
    When I do GET /api/v4/cat/job-executions/test-job-instruction-cancel-job-execution-running
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "status": 0
    }
    """
    When I do PUT /api/v4/cat/jobs/test-job-instruction-cancel-job-execution-running:
    """
    {
      "name": "test-job-instruction-cancel-job-execution-running",
      "config": "test-job-config-running",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/job-executions/test-job-instruction-cancel-job-execution-running
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "status": 3
    }
    """

  Scenario: given instruction should cancel running executions on instruction update
    When I am admin
    When I do GET /api/v4/cat/job-executions/test-job-instruction-cancel-execution-running-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "status": 0
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-cancel-execution-running-2:
    """
    {
      "name": "test-instruction-cancel-execution-running-2-name",
      "entity_patterns": [
        {
          "name": "test filter"
        }
      ],
      "description": "test-instruction-cancel-execution-running-2-description",
      "enabled": true,
      "steps": [
        {
          "name": "test-instruction-cancel-execution-running-2-step-1",
          "operations": [
            {
              "name": "test-instruction-cancel-execution-running-2-step-1-operation-1",
              "time_to_complete": {"seconds": 1, "unit":"s"},
              "description": "test-instruction-cancel-execution-running-2-step-1-operation-1-description",
              "jobs": ["test-job-instruction-cancel-execution-running-2"]
            },
            {
              "name": "test-instruction-cancel-execution-running-2-step-1-operation-2",
              "time_to_complete": {"seconds": 3, "unit":"s"},
              "description": "test-instruction-cancel-execution-running-2-step-1-operation-2-description"
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-cancel-execution-running-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/job-executions/test-job-instruction-cancel-execution-running-2
    Then the response code should be 410