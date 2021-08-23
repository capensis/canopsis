Feature: get a running job
  I need to be able to get a running job
  Only admin should be able to get a running job

  Scenario: given job should start job for operation of instruction
    When I am admin
    When I do GET /api/v4/cat/job-executions/test-job-execution-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-execution-to-get",
      "job_id": "test-job-execution-to-get",
      "name": "test-job-execution-to-get",
      "status": 2,
      "fail_reason": "test-job-execution-to-get-fail-reason",
      "output": "test-job-execution-to-get-output",
      "payload": "",
      "query": null,
      "started_at": 1597906527,
      "launched_at": 1597906527,
      "completed_at": 1597906537
    }
    """

  Scenario: given unauth request should not allow access
    When I do GET /api/v4/cat/job-executions/test
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/job-executions/test
    Then the response code should be 403
