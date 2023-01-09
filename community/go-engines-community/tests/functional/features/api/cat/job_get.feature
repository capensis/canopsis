Feature: get a job
  I need to be able to get a job
  Only admin should be able to get a job

  Scenario: given get all request should return jobs
    When I am admin
    When I do GET /api/v4/cat/jobs?search=test-job-name-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-job-to-get-1",
          "name": "test-job-name-to-get-1",
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
        },
        {
          "_id": "test-job-to-get-2",
          "name": "test-job-name-to-get-2",
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
          "payload": "",
          "query": {
            "key1": "val1",
            "key2": "val2"
          },
          "multiple_executions": true,
          "retry_amount": 3,
          "retry_interval": {
            "value": 10,
            "unit": "s"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get all request should return jobs with flags
    When I am admin
    When I do GET /api/v4/cat/jobs?search=test-job-name-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-to-get-1",
          "deletable": true,
          "running": false
        },
        {
          "_id": "test-job-to-get-2",
          "deletable": true,
          "running": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/jobs
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/jobs
    Then the response code should be 403

  Scenario: given get request should return ok
    When I am admin
    When I do GET /api/v4/cat/jobs/test-job-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-to-get-1",
      "name": "test-job-name-to-get-1",
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

  Scenario: given job with linked instruction should return corresponding flags
    When I am admin
    When I do GET /api/v4/cat/jobs?search=test-job-to-check-linked-to-manual-instruction&with_flags=true&sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-to-check-linked-to-manual-instruction",
          "deletable": false,
          "running": false
        },
        {
          "_id": "test-job-to-check-linked-to-manual-instruction-execution",
          "deletable": false,
          "running": false
        }
      ]
    }
    """
    When I do GET /api/v4/cat/jobs?search=test-job-to-check-linked-to-auto-instruction&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-to-check-linked-to-auto-instruction",
          "deletable": false,
          "running": false
        }
      ]
    }
    """

  Scenario: given get request with not exist job should return not found error
    When I am admin
    When I do GET /api/v4/cat/jobs/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/jobs/test-job-to-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/jobs/test-job-to-get
    Then the response code should be 403
