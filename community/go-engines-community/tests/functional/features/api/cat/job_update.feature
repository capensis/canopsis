Feature: Job update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "author": "new-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
    }
    """
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "author": "new-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
    }
    """
    Then the response code should be 403

  Scenario: PUT a valid job without any changes
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-to-update",
      "author": "test-author",
      "config": {
        "_id": "test-job-config-to-link",
        "name": "test-job-config-name-to-link",
        "type": "rundeck",
        "host": "http://example.com",
        "author": "test-author",
        "auth_token": "test-auth-token"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """

  Scenario: PUT a valid job
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "author": "new-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-to-update",
      "author": "new-author",
      "config": {
        "_id": "test-job-config-to-link",
        "name": "test-job-config-name-to-link",
        "type": "rundeck",
        "host": "http://example.com",
        "author": "test-author",
        "auth_token": "test-auth-token"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """

  Scenario: PUT a valid job that doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update-do-not-exists:
    """
    {
      "name": "test-job-name-do-not-exists",
      "author": "new-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 404

  Scenario: PUT an invalid job where name already exists
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-get",
      "author": "new-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name already exists"
      }
    }
    """

  Scenario: PUT an invalid job where config doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "author": "new-author",
      "config": "test-job-config-not-exist",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "job's config doesn't exist"
    }
    """