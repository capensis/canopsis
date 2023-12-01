Feature: Get icons
  I need to be able to get icons
  Only admin should be able to get icons

  Scenario: given search request should return ok
    When I am admin
    When I add form field title=test-icon-to-get-1-1-title
    And I add form file file=test.svg
    And I do POST /api/v4/icons
    Then the response code should be 201
    When I save response iconId1={{ .lastResponse._id }}
    When I add form field title=test-icon-to-get-1-2-title
    And I add form file file=test2.svg
    And I do POST /api/v4/icons
    Then the response code should be 201
    When I save response iconId2={{ .lastResponse._id }}
    When I am unauth
    When I do GET /api/v4/icons?search=test-icon-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "test-icon-to-get-1-1-title"
        },
        {
          "title": "test-icon-to-get-1-2-title"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do GET /api/v4/icons/{{ .iconId1 }}
    Then the response code should be 200
    Then the response body should contain file test.svg
    When I do GET /api/v4/icons/{{ .iconId2 }}
    Then the response code should be 200
    Then the response body should contain file test2.svg
