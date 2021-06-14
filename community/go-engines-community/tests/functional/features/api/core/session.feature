Feature: Session auth user
  I need to be able to auth user by session.

  Scenario: given user username and password should log in user and create session
    When I do POST /auth:
    """
      {
        "username": "root",
        "password": "test"
      }
    """
    Then the response code should be 200
    Then the response body should be:
    """
	  {
		"crecord_name": "root",
		"authkey": "2be1d0ee-5e9e-4b1f-9d80-e857e50f4191",
		"role": "admin",
		"contact": {
		  "name": "John Doe",
          "address": "30 Rue du Triez 59290 Wasquehal France"
        },
		"mail": "admin@canopsis.net"
	  }
	"""

  Scenario: given auth user by session cookie should allow access
    When I do POST /auth:
    """
      {
        "username": "root",
        "password": "test"
      }
    """
    When I do GET /api/v4/alarms
    Then the response code should be 200

  Scenario: given auth user by session cookie should clear sessions
    When I do POST /auth:
    """
      {
        "username": "root",
        "password": "test"
      }
    """
    When I do GET /logout
    Then the response code should be 200

  Scenario: given no session should not allow access
    When I do POST /auth:
    """
      {
        "username": "root",
        "password": "test"
      }
    """
    When I do GET /logout
    When I do GET /api/v4/alarms
    Then the response code should be 401

  Scenario: given invalid username and password should return error
    When I do POST /auth:
    """
      {
        "username": "nouser",
        "password": "nopass"
      }
    """
    Then the response code should be 401

  Scenario: given no username and password should return errors
    When I do POST /auth
    Then the response code should be 400
    Then the response body should be:
    """
      {
        "errors": {
          "password":"Password is missing.",
          "username":"Username is missing."
        }
      }
	"""