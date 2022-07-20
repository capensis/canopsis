Feature: create a job
  I need to be able to create a job
  Only admin should be able to create a job

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/jobs
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/jobs
    Then the response code should be 403

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "name": "test-job-to-create-1",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"resource\": \"{{ `{{ .Alarm.Value.Resource }}` }}\",\"entity\": \"{{ `{{ .Entity.ID }}` }}\"}",
      "query": {
        "resource": "{{ `{{ .Alarm.Value.Resource }}` }}"
      },
      "multiple_executions": true,
      "retry_amount": 3,
      "retry_interval": {
        "value": 10,
        "unit": "s"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-create-1",
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
      "multiple_executions": true,
      "retry_amount": 3,
      "retry_interval": {
        "value": 10,
        "unit": "s"
      }
    }
    """
    When I do GET /api/v4/cat/jobs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-job-to-create-1",
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
      "multiple_executions": true,
      "retry_amount": 3,
      "retry_interval": {
        "value": 10,
        "unit": "s"
      }
    }
    """

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "_id": "test-job-to-create-2",
      "name": "test-job-to-create-2-name",
      "config": "test-job-config-to-edit-job",
      "job_id": "test-job-id",
      "payload": "{\"key1\": \"val1\",\"key2\": \"val2\"}",
      "multiple_executions": false
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/jobs/test-job-to-create-2
    Then the response code should be 200

  Scenario: given create request with custom id that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/cat/jobs:
    """json
    {
      "_id": "test-job-to-check-unique",
      "name": "test-job-to-create-3-name",
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

  Scenario: given create request the name that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/cat/jobs:
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

  Scenario: given invalid create request should return errors
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
