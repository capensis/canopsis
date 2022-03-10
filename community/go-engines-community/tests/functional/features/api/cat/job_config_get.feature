Feature: get a job's config
  I need to be able to get a job's config
  Only admin should be able to get a job's config

  Scenario: given get all request should return job's configs
    When I am admin
    When I do GET /api/v4/cat/job-configs?search=test-job-config-name-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-job-config-to-get",
          "auth_token": "test-auth-token",
          "auth_username": "",
          "host": "http://example.com",
          "name": "test-job-config-name-to-get",
          "type": "rundeck",
          "author": {
            "_id": "root",
            "name": "root"
          }
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

  Scenario: given get all request should return job's configs with flags
    When I am admin
    When I do GET /api/v4/cat/job-configs?search=test-job-config-name-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-config-to-get",
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

  Scenario: GET a job's config but unauthorized
    When I do GET /api/v4/cat/job-configs/test-job-config-to-get
    Then the response code should be 401

  Scenario: GET a job's config but without permissions
    When I am noperms
    When I do GET /api/v4/cat/job-configs/test-job-config-to-get
    Then the response code should be 403

  Scenario: Get a job's config with success
    When I am admin
    When I do GET /api/v4/cat/job-configs/test-job-config-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-config-to-get",
      "auth_token": "test-auth-token",
      "auth_username": "",
      "host": "http://example.com",
      "name": "test-job-config-name-to-get",
      "type": "rundeck",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  Scenario: Get a job's config with linked job
    When I am admin
    When I do GET /api/v4/cat/job-configs?search=test-job-config-to-check-linked&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-job-config-to-check-linked",
          "deletable": false,
          "running": false
        }
      ]
    }
    """

  Scenario: Get a job's config with not found response
    When I am admin
    When I do GET /api/v4/cat/job-configs/test-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
