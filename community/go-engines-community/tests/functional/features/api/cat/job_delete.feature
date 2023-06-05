Feature: delete a job

  Scenario: given delete unauth request should not allow access
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 401

  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 403

  Scenario: given delete request should delete job
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 204

  Scenario: given not exist job should return error
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given linked job should return error
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-to-check-linked-to-manual-instruction
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "job is linked with instruction"
    }
    """
    When I do DELETE /api/v4/cat/jobs/test-job-to-check-linked-to-auto-instruction
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "job is linked with instruction"
    }
    """
    When I do DELETE /api/v4/cat/jobs/test-job-to-check-linked-to-manual-instruction-execution
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "job is linked with instruction"
    }
    """
