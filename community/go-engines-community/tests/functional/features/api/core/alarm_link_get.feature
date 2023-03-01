Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get rule with templates should return links
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-to-alarm-link-get-2?ids[]=test-alarm-to-link-get-1&ids[]=test-alarm-to-link-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "label": "test-link-rule-to-alarm-link-get-2-link-1-label",
        "icon_name": "test-link-rule-to-alarm-link-get-2-link-1-icon",
        "url": "http://test-link-rule-to-alarm-link-get-2-link-1-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val&info[]=test-resource-to-alarm-link-get-2-info-1-val&"
      },
      {
        "label": "test-link-rule-to-alarm-link-get-2-link-2-label",
        "icon_name": "test-link-rule-to-alarm-link-get-2-link-2-icon",
        "url": "http://test-link-rule-to-alarm-link-get-2-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val&info[]=test-resource-to-alarm-link-get-2-info-1-val&"
      }
    ]
    """

  @concurrent
  Scenario: given get rule with source code should return links
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-to-alarm-link-get-4?ids[]=test-alarm-to-link-get-1&ids[]=test-alarm-to-link-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    [
      {
        "label": "test-link-rule-to-alarm-link-get-4-link-1-label",
        "icon_name": "test-link-rule-to-alarm-link-get-4-link-1-icon",
        "url": "http://test-link-rule-to-alarm-link-get-4-link-1-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
      },
      {
        "label": "test-link-rule-to-alarm-link-get-4-link-2-label",
        "icon_name": "test-link-rule-to-alarm-link-get-4-link-2-icon",
        "url": "http://test-link-rule-to-alarm-link-get-4-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
      }
    ]
    """

  @concurrent
  Scenario: given get links unauth request should not allow access
    When I do GET /api/v4/alarm-links/test-link-rule-not-exist
    Then the response code should be 401

  @concurrent
  Scenario: given get links request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarm-links/test-link-rule-not-exist
    Then the response code should be 403

  @concurrent
  Scenario: given get links invalid request should return errors
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-to-alarm-link-get-1
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "ids": "Ids is missing."
      }
    }
    """

  @concurrent
  Scenario: given not exist rule should return errors
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-not-exist?ids[]=test-alarm-to-link-get-1&ids[]=test-alarm-to-link-get-2
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given not matched alarm should return errors
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-to-alarm-link-get-1?ids[]=test-alarm-to-link-get-1&ids[]=test-alarm-to-link-get-2
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "ids": "Alarms aren't matched to rule."
      }
    }
    """

  @concurrent
  Scenario: given not exist alarm should return errors
    When I am admin
    When I do GET /api/v4/alarm-links/test-link-rule-to-alarm-link-get-1?ids[]=test-alarm-not-exist
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "ids": "Alarms aren't matched to rule."
      }
    }
    """
