Feature: delete a job

  Scenario: DELETE a job but unauthorized
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 401

  Scenario: DELETE a job but without permissions
    When I am noperms
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 403

  Scenario: DELETE a job with success
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-to-delete
    Then the response code should be 204

  Scenario: DELETE a job with not found response
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
  Scenario: DELETE a linked job shouldn't be possible
    When I am admin
    When I do DELETE /api/v4/cat/jobs/test-job-to-test-instruction-to-get-step-2-operation-1-2
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "job is linked with instruction"
    }
    """