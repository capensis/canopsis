Feature: Get entities by pbehavior id
  I need to be able to get entities by pbehavior id
  Only admin should be able to get entities by pbehavior id

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/pbehaviors/test-not-found/entities
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/pbehaviors/test-not-found/entities
    Then the response code should be 403
    
  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/pbehaviors/test-not-found/entities
    Then the response code should be 404

  Scenario: given get request should return entities
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-1/entities
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-1/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-2/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-3/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-4/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-5/test-component-default"
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
    
  Scenario: given paginated get request should return entities by page
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-1/entities?page=1&limit=2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-1/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-2/test-component-default"
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
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-1/entities?page=2&limit=2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-3/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-4/test-component-default"
        }
      ],
      "meta": {
        "page": 2,
        "per_page": 2,
        "page_count": 3,
        "total_count": 5
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-1/entities?page=3&limit=2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-5/test-component-default"
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

  Scenario: given sorted get request should return entities by sort
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-1/entities?sort_by=name&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-5/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-4/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-3/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-2/test-component-default"
        },
        {
          "_id": "test-resource-to-entities-get-by-pbh-1/test-component-default"
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

  Scenario: given get request by pbehavior with old mongo query should return entities
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-entities-get-2/entities
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-to-entities-get-by-pbh-6/test-component-default"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
