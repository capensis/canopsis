Feature: Job update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/jobs/test-job-to-update
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/cat/jobs/test-job-to-update
    Then the response code should be 403

  Scenario: PUT a valid job without any changes
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """json
    {
      "name": "test-job-name-to-update",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-to-update",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "config": {
        "_id": "test-job-config-to-edit-job",
        "name": "test-job-config-to-edit-job-name",
        "type": "rundeck",
        "host": "http://example.com",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "auth_username": "",
        "auth_token": "test-auth-token"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "query": null,
      "multiple_executions": false
    }
    """

  Scenario: PUT a valid job
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """json
    {
      "name": "test-job-name-to-update",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-to-update",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "config": {
        "_id": "test-job-config-to-edit-job",
        "name": "test-job-config-to-edit-job-name",
        "type": "rundeck",
        "host": "http://example.com",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "auth_username": "",
        "auth_token": "test-auth-token"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "query": null,
      "multiple_executions": true
    }
    """

  Scenario: PUT a valid job that doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update-do-not-exists:
    """json
    {
      "name": "test-job-name-do-not-exists",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 404

  Scenario: PUT an invalid job where name already exists
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """json
    {
      "name": "test-job-to-check-unique-name",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  Scenario: PUT an invalid job where config doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/jobs/test-job-to-update:
    """json
    {
      "name": "test-job-name-to-update",
      "config": "test-job-config-not-exist",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "config": "Config doesn't exist."
      }
    }
    """
