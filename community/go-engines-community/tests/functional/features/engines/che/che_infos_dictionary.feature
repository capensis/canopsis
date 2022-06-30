Feature: get infos dictionary
  I need to be able to get infos dictionary

  Scenario: given resource check event and truncate output and long_output field
    Given I am admin
    When I wait the next periodical process
    When I do GET /api/v4/entity-infos-dictionary?search=test-entity-infos-dictionary-key
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-key-1"
        },
        {
          "value": "test-entity-infos-dictionary-key-2"
        },
        {
          "value": "test-entity-infos-dictionary-key-3"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/entity-infos-dictionary?search=test-entity-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-key-2"
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
    When I do GET /api/v4/entity-infos-dictionary?key=test-entity-infos-dictionary-key-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-value-11"
        },
        {
          "value": "test-entity-infos-dictionary-value-2"
        },
        {
          "value": "test-entity-infos-dictionary-value-5"
        },
        {
          "value": "test-entity-infos-dictionary-value-8"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 4
      }
    }
    """
    When I do GET /api/v4/entity-infos-dictionary?key=test-entity-infos-dictionary-key-2&search=test-entity-infos-dictionary-value-5
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "value": "test-entity-infos-dictionary-value-5"
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
    When I do GET /api/v4/entity-infos-dictionary?key=test-entity-infos-dictionary-should-be-ignored
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
