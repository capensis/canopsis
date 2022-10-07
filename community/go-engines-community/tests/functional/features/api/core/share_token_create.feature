Feature: Create a share token
  I need to be able to create a share token
  Only admin should be able to create a share token

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {
      "description": "test-share-token-to-create-1",
      "duration": {
        "value": 7,
        "unit": "d"
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "user": {
        "_id": "root",
        "name": "root"
      },
      "role": {
        "_id": "admin",
        "name": "admin"
      },
      "description": "test-share-token-to-create-1"
    }
    """
    When I save response now={{ now }}
    When I save response expectedExpired={{ nowAdd "7d" }}
    When I save response expired={{ .lastResponse.expired }}
    When I save response created={{ .lastResponse.created }}
    Then the difference between now created is in range -2,2
    Then the difference between expectedExpired expired is in range -2,2
    When I do GET /api/v4/share-tokens?search={{ .lastResponse.value }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "user": {
            "_id": "root",
            "name": "root"
          },
          "role": {
            "_id": "admin",
            "name": "admin"
          },
          "description": "test-share-token-to-create-1"
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

  Scenario: given create request without expiration should return ok
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {}
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "user": {
        "_id": "root",
        "name": "root"
      },
      "role": {
        "_id": "admin",
        "name": "admin"
      },
      "description": "",
      "expired": null
    }
    """
    When I do GET /api/v4/share-tokens?search={{ .lastResponse.value }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "user": {
            "_id": "root",
            "name": "root"
          },
          "role": {
            "_id": "admin",
            "name": "admin"
          },
          "description": "",
          "expired": null
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

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/share-tokens
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/share-tokens
    Then the response code should be 403
