Feature: create a job's config
  I need to be able to create a job's config
  Only admin should be able to create a job's config

  Scenario: POST a valid job's config but unauthorized
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 401

  Scenario: POST a valid job's config but without permissions
    When I am noperms
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 403

  Scenario: POST a valid job's config
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-job-config-name",
      "host": "http://example.com",
      "type": "rundeck",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "auth_token": "test-auth-token"
    }
    """

  Scenario: POST a valid job's config
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name-2",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/job-configs/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: POST a valid job's config with custom id
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "_id": "custom-id",
      "name": "test-job-config-name-2-custom-id-1",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/job-configs/custom-id
    Then the response code should be 200

  Scenario: POST a valid job's config with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "_id": "test-job-config-to-update",
      "name": "test-job-config-name-2-custom-id-2",
      "host": "http://example.com",
      "type": "rundeck",
      "auth_token": "test-auth-token"
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

  Scenario: POST an invalid job's config, where required fields are empty
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "host": "Host is missing.",
        "name": "Name is missing.",
        "type": "Type is missing."
      }
    }
    """
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "type": "rundeck"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "auth_token": "AuthToken is missing."
      }
    }
    """
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "type": "awx"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "auth_token": "AuthToken is missing."
      }
    }
    """
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "type": "jenkins"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "auth_token": "AuthToken is missing.",
        "auth_username": "AuthUsername is missing."
      }
    }
    """

  Scenario: POST an invalid job's config, where host is not an url
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name-3",
      "host": "abc",
      "type": "rundeck",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "host": "abc is not an url."
      }
    }
    """

  Scenario: POST an invalid job's config, where type is not valid
    When I am admin
    When I do POST /api/v4/cat/job-configs:
    """json
    {
      "name": "test-job-config-name-3",
      "host": "http://example.com",
      "type": "not-valid",
      "auth_token": "test-auth-token"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type must be one of [awx, jenkins, rundeck]."
      }
    }
    """
