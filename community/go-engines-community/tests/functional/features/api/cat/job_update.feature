Feature: Job update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
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
      "author": {"_id": "root", "name": "root"},
      "config": {
        "_id": "test-job-config-to-link",
        "name": "test-job-config-name-to-link",
        "type": "rundeck",
        "host": "http://example.com",
        "author": {"_id": "test-author", "name": "test-author-username"},
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
      "author": {"_id": "root", "name": "root"},
      "config": {
        "_id": "test-job-config-to-link",
        "name": "test-job-config-name-to-link",
        "type": "rundeck",
        "host": "http://example.com",
        "author": {"_id": "test-author", "name": "test-author-username"},
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
        "name": "Name already exists."
      }
    }
    """

  Scenario: PUT an invalid job where config doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """
    {
      "name": "test-job-name-to-update",
      "config": "test-job-config-not-exist",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {"config": "Config doesn't exist."}
    }
    """