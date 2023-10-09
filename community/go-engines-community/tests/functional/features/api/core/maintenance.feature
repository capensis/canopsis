Feature: maintenance feature
  I need to be able to set canopsis to maintenance mode
  Only admin should be able to set canopsis to maintenance mode

  @standalone
  Scenario: given request with various users, only admin user should be able to enable/disable canopsis maintenance mode
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "color": "#BBBBBB",
      "message": "test maintenance"
    }
    """
    Then the response code should be 401
    When I am noperms
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "color": "#BBBBBB",
      "message": "test maintenance"
    }
    """
    Then the response code should be 403
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "color": "#BBBBBB",
      "message": "test maintenance"
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "maintenance": true
    }
    """
    When I do GET /api/v4/active-broadcast-message
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "message": "test maintenance",
        "maintenance": true,
        "color": "#BBBBBB"
      },
      {}
    ]
    """
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "maintenance": false
    }
    """
    When I do GET /api/v4/active-broadcast-message
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {}
    ]
    """

  @standalone
  Scenario: given enable/disable maintenance request, should be enabled/disabled only once, enable request should contain message field
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "message": "message is required"
      }
    }
    """
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "message": "test maintenance",
      "color": "#BBBBBB"
    }
    """
    Then the response code should be 204
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "message": "test maintenance",
      "color": "#BBBBBB"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "maintenance mode has already been enabled"
    }
    """
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "maintenance mode has already been disabled"
    }
    """

  @standalone
  Scenario: given users and maintenance mode, only admin user should be able to login or be authenticated,
    others should be able to login again after maintenance is disabled
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "message": "test maintenance",
      "color": "#BBBBBB"
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/login:
    """json
    {
      "username": "manageruser",
      "password": "test"
    }
    """
    Then the response code should be 503
    Then the response body should be:
    """json
    {
      "error": "canopsis is under maintenance"
    }
    """
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/app-info
    Then the response code should be 503
    Then the response body should be:
    """json
    {
      "error": "canopsis is under maintenance"
    }
    """
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I do POST /api/v4/login:
    """json
    {
      "username": "manageruser",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I am authenticated with username "manageruser" and password "test"
    When I do GET /api/v4/app-info
    Then the response code should be 200

  @standalone
  Scenario: given user with token, when maintenance mode is enabled it should remove user's token
    When I do POST /api/v4/login:
    """json
    {
      "username": "manageruser",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I save response token={{ .lastResponse.access_token }}
    When I set header Authorization=Bearer {{ .token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "message": "test maintenance",
      "color": "#BBBBBB"
    }
    """
    Then the response code should be 204
    When I set header Authorization=Bearer {{ .token }}
    When I do GET /api/v4/app-info
    Then the response code should be 401
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I set header Authorization=Bearer {{ .token }}
    When I do GET /api/v4/app-info
    Then the response code should be 401

  @standalone
  Scenario: given user with admin/not admin shared tokens, when maintenance mode is enabled it shouldn't remove shared token,
    but should restrict access for not admin shared tokens while maintenance is active
    When I am manager
    When I do POST /api/v4/share-tokens:
    """json
    {
      "description": "test",
      "duration": {
        "value": 7,
        "unit": "d"
      }
    }
    """
    Then the response code should be 201
    When I save response manager_shared_token={{ .lastResponse.value }}
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {
      "description": "test",
      "duration": {
        "value": 7,
        "unit": "d"
      }
    }
    """
    Then the response code should be 201
    When I save response admin_shared_token={{ .lastResponse.value }}
    When I set header Authorization=Bearer {{ .manager_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I set header Authorization=Bearer {{ .admin_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": true,
      "message": "test maintenance",
      "color": "#BBBBBB"
    }
    """
    Then the response code should be 204
    When I set header Authorization=Bearer {{ .manager_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 503
    Then the response body should be:
    """json
    {
      "error": "canopsis is under maintenance"
    }
    """
    When I set header Authorization=Bearer {{ .admin_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I am admin
    When I do PUT /api/v4/maintenance:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I set header Authorization=Bearer {{ .manager_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200
    When I set header Authorization=Bearer {{ .admin_shared_token }}
    When I do GET /api/v4/app-info
    Then the response code should be 200

