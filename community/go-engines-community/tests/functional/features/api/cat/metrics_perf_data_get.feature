Feature: Get perf data metrics
  I need to be able to get perf data metrics
  Only admin should be able to get perf data metrics

  @concurrent
  Scenario: given get request should return metrics
    When I am admin
    When I do GET /api/v4/cat/perf-data-metrics?search=test-perf-data-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        "test-perf-data-to-get-1",
        "test-perf-data-to-get-2"
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/perf-data-metrics
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/perf-data-metrics
    Then the response code should be 403
