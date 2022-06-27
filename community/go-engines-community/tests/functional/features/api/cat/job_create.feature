Feature: create a job
  I need to be able to create a job
  Only admin should be able to create a job

  Scenario: POST a valid job but unauthorized
    When I do POST /api/v4/cat/jobs
    Then the response code should be 401

  Scenario: POST a valid job but without permissions
    When I am noperms
    When I do POST /api/v4/cat/jobs
    Then the response code should be 403

  Scenario: POST a valid job
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-name",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "query": {
        "resource": "{{ `{{ .Alarm.Value.Resource }}` }}"
      },
      "multiple_executions": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-job-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "config": {
        "_id": "test-job-config-to-edit-job",
        "auth_token": "test-auth-token",
        "host": "http://example.com",
        "name": "test-job-config-to-edit-job-name",
        "type": "rundeck"
      },
      "job_id": "test-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "query": {
        "resource": "{{ `{{ .Alarm.Value.Resource }}` }}"
      },
      "multiple_executions": true
    }
    """

  Scenario: POST a valid job
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-name-2",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/jobs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-job-name-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "config": {
        "_id": "test-job-config-to-edit-job",
        "auth_token": "test-auth-token",
        "host": "http://example.com",
        "name": "test-job-config-to-edit-job-name",
        "type": "rundeck"
      },
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """

  Scenario: POST a valid job with custom id
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "_id": "custom-id",
      "name": "test-job-name-2-custom-id-1",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/jobs/custom-id
    Then the response code should be 200

  Scenario: POST a valid job with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "_id": "test-job-to-update",
      "name": "test-job-name-2-custom-id-2",
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
        "_id": "ID already exists."
      }
    }
    """

  Scenario: POST an invalid job where name already exists
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-name-to-get",
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

  Scenario: POST an invalid job where config doesn't exist
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-name-with-not-existed-config",
      "config": "test-job-config-not-exists",
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
