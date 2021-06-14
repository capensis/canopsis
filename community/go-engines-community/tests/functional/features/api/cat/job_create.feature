Feature: create a job
  I need to be able to create a job
  Only admin should be able to create a job

  Scenario: POST a valid job but unauthorized
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 401

  Scenario: POST a valid job but without permissions
    When I am noperms
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 403

  Scenario: POST a valid job
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "name": "test-job-name",
      "author": "test-author",
      "config": {
        "_id": "test-job-config-to-link",
        "auth_token": "test-auth-token",
        "host": "http://example.com",
        "name": "test-job-config-name-to-link",
        "type": "rundeck"
      },
      "job_id": "test-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}"
    }
    """

  Scenario: POST a valid job
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name-2",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/jobs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-job-name-2",
      "author": "test-author",
      "config": {
        "_id": "test-job-config-to-link",
        "auth_token": "test-auth-token",
        "host": "http://example.com",
        "name": "test-job-config-name-to-link",
        "type": "rundeck"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """

  Scenario: POST a valid job with custom id
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "_id": "custom-id",
      "name": "test-job-name-2-custom-id-1",
      "author": "test-author",
      "config": "test-job-config-to-link",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/jobs/custom-id
    Then the response code should be 200

  Scenario: POST a valid job with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-2-custom-id-2",
      "author": "test-author",
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
        "_id": "ID already exists"
      }
    }
    """

  Scenario: POST an invalid job where name already exists
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name-to-get",
      "author": "test-author",
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

  Scenario: POST an invalid job where config doesn't exist
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """
    {
      "name": "test-job-name-with-not-existed-config",
      "author": "test-author",
      "config": "test-job-config-not-exists",
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