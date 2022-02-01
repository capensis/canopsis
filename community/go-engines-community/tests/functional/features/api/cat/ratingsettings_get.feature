Feature: Get a rating settings
  I need to be able to get a rating settings
  Only admin should be able to get a rating settings

  Scenario: given search request should return rating settings
    When I am admin
    When I do GET /api/v4/cat/rating-settings?search=test-rating-settings-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "id": 4,
          "label": "test-rating-settings-to-get-1",
          "enabled": true
        },
        {
          "id": 5,
          "label": "infos.test-rating-settings-to-get-2",
          "enabled": false
        },
        {
          "id": 6,
          "label": "test-rating-settings-to-get-3",
          "enabled": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given paginated request should return rating settings
    When I am admin
    When I do GET /api/v4/cat/rating-settings?page=1&limit=2&search=test-rating-settings-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "id": 4,
          "label": "test-rating-settings-to-get-1",
          "enabled": true
        },
        {
          "id": 5,
          "label": "infos.test-rating-settings-to-get-2",
          "enabled": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 2,
        "per_page": 2,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/cat/rating-settings?page=2&limit=2&search=test-rating-settings-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "id": 6,
          "label": "test-rating-settings-to-get-3",
          "enabled": false
        }
      ],
      "meta": {
        "page": 2,
        "page_count": 2,
        "per_page": 2,
        "total_count": 3
      }
    }
    """

  Scenario: given enabled request should return rating settings
    When I am admin
    When I do GET /api/v4/cat/rating-settings?enabled=true&search=test-rating-settings-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "id": 4,
          "label": "test-rating-settings-to-get-1",
          "enabled": true
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

  Scenario: given search request should return nothing
    When I am admin
    When I do GET /api/v4/cat/rating-settings?search=not-exist
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/rating-settings
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/rating-settings
    Then the response code should be 403
