Feature: Account auth user
  I need to be able to get auth user account.

  Scenario: given user username and password should get account
    When I do POST /auth:
    """
    {
      "username": "test-user-to-get",
      "password": "test"
    }
    """
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-user-to-get",
      "authkey": "4de736d3-b2e9-451c-9367-bcf15acc3309",
      "crecord_name": "test-user-to-get",
      "firstname": "John",
      "lastname": "Doe",
      "enable": true,
      "mail": "test-user-to-get@canopsis.net",
      "external_id": "",
      "source": "",
      "ui_language": "fr",
      "groupsNavigationType": "side-bar",
      "tours": {"alarmsExpandPanasdel": true},
      "defaultview": "a1d329c8-9939-40bc-bb86-b63172ac44b4",
      "role": "manager",
      "rights": {
        "api_action": {
          "checksum": 15
        },
        "api_alarm_read": {
          "checksum": 1
        },
        "api_entity_read": {
          "checksum": 1
        },
        "api_event": {
          "checksum": 1
        },
        "api_execution": {
          "checksum": 1
        },
        "api_heartbeat": {
          "checksum": 15
        },
        "api_instruction": {
          "checksum": 15
        },
        "api_job": {
          "checksum": 15
        },
        "api_job_config": {
          "checksum": 15
        },
        "api_pbehavior": {
          "checksum": 15
        },
        "api_pbehaviorexception": {
          "checksum": 15
        },
        "api_pbehaviorreason": {
          "checksum": 15
        },
        "api_pbehaviortype": {
          "checksum": 15
        },
        "api_watcher": {
          "checksum": 15
        },
        "api_webhook": {
          "checksum": 15
        }
      }
    }
	"""

  Scenario: given unauth request should not allow access
    When I do GET /api/v4/account/me
    Then the response code should be 401
