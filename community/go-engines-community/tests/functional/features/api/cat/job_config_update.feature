Feature: Job's config update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update
    Then the response code should be 403

  Scenario: PUT a valid job's config
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update:
    """json
    {
      "name": "test-job-config-name-to-update-new",
      "type": "rundeck",
      "host": "http://example-2.com",
      "auth_token": "new token"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-config-to-update",
      "auth_token": "new token",
      "auth_username": "",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "host": "http://example-2.com",
      "name": "test-job-config-name-to-update-new",
      "type": "rundeck"
    }
    """

  Scenario: PUT a valid job's config that doesn't exist
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update-do-not-exists:
    """json
    {
      "name": "test-job-config-name-do-not-exists",
      "type": "rundeck",
      "host": "http://example.com",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 404

  Scenario: PUT a valid job's config with invalid host
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update:
    """json
    {
      "name": "test-job-config-name-to-update",
      "type": "rundeck",
      "host": "abc",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
      {
        "errors": {
          "host": "abc is not an url."
        }
      }
    """

  Scenario: PUT a valid job's config without any changes
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update:
    """json
    {
      "auth_token": "new token",
      "host": "http://example-2.com",
      "name": "test-job-config-name-to-update-new",
      "type": "rundeck"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-job-config-to-update",
      "auth_token": "new token",
      "auth_username": "",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "host": "http://example-2.com",
      "name": "test-job-config-name-to-update-new",
      "type": "rundeck"
    }
    """

  Scenario: PUT an invalid job's config with wrong type
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update:
    """json
    {
      "auth_token": "new token",
      "host": "http://example-2.com",
      "name": "test-job-config-name-to-update-new",
      "type": "wrong_type"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
      {
        "errors": {
          "type": "Type must be one of [awx, jenkins, rundeck]."
        }
      }
    """

  Scenario: PUT an invalid job's config with already existed name
    When I am admin
    When I do PUT /api/v4/cat/job-configs/test-job-config-to-update:
    """json
    {
      "auth_token": "new token",
      "host": "http://example-2.com",
      "name": "test-job-config-name-to-get",
      "type": "rundeck"
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
