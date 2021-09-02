Feature: delete a job's config

  Scenario: DELETE a job's config but unauthorized
    When I do DELETE /api/v4/cat/job-configs/test-job-config-to-delete
    Then the response code should be 401

  Scenario: DELETE a job's config but without permissions
    When I am noperms
    When I do DELETE /api/v4/cat/job-configs/test-job-config-to-delete
    Then the response code should be 403

  Scenario: DELETE a job's config with success
    When I am admin
    When I do DELETE /api/v4/cat/job-configs/test-job-config-to-delete
    Then the response code should be 204

  Scenario: DELETE a job's config with not found response
    When I am admin
    When I do DELETE /api/v4/cat/job-configs/test-job-config-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
  Scenario: DELETE a linked job's config shouldn't be possible
    When I am admin
    When I do DELETE /api/v4/cat/job-configs/test-job-config-to-check-linked
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "job's config is linked with job"
    }
    """