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
          "_id": "test-job-to-get",
          "name": "test-job-name-to-get",
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
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
          "_id": "test-job-to-get",
          "deletable": true,
          "running": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: GET a job but unauthorized
    When I do GET /api/v4/cat/jobs/test-job-to-get
    Then the response code should be 401

  Scenario: GET a job but without permissions
    When I am noperms
    When I do GET /api/v4/cat/jobs/test-job-to-get
    Then the response code should be 403

  Scenario: Get a job with success
    When I am admin
    When I do GET /api/v4/cat/jobs/test-job-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-to-get",
      "name": "test-job-name-to-get",
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

  Scenario: Get a job with linked instruction
    When I am admin
    When I do GET /api/v4/cat/jobs?search=test-job-to-check-linked-to-manual-instruction&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-to-check-linked-to-manual-instruction",
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

  Scenario: Get a job with not found response
    When I am admin
    When I do GET /api/v4/cat/jobs/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
