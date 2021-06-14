Feature: get PBehavior eids
  I need to be able to get PBehavior eids
  Only admin should be able to get PBehavior eids

  Scenario: GET unauthorized
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids
    Then the response code should be 403

  Scenario: GET success
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-1/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-2/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-3/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-4/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-5/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 5
      }
    }
    """
  Scenario: GET success with pagination
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-1/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-2/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-3/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-4/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-5/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 5
      }
    }
    """
  Scenario: GET success with pagination
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids?page=1&limit=2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-1/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-2/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 2,
        "page_count": 3,
        "total_count": 5
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids?page=3&limit=2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-5/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 3,
        "per_page": 2,
        "page_count": 3,
        "total_count": 5
      }
    }
    """
  Scenario: GET success with sort
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids?page=1&limit=2&sort_by=id&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-5/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-4/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 2,
        "page_count": 3,
        "total_count": 5
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-get-eids/eids?page=1&limit=2&sort_by=id&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "id": "test-alarm-get-pbehavior-eids-1/test-alarm-get-pbehavior-eids-component"
        },
        {
          "id": "test-alarm-get-pbehavior-eids-2/test-alarm-get-pbehavior-eids-component"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 2,
        "page_count": 3,
        "total_count": 5
      }
    }
    """