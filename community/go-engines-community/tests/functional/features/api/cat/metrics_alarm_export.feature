Feature: Export alarm metrics
  I need to be able to export alarm metrics
  Only admin should be able to export alarm metrics

  Scenario: given export request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_alarms&sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    total_alarms,{{ nowDateAdd "-72h" }},0
    total_alarms,{{ nowDateAdd "-48h" }},1
    total_alarms,{{ nowDateAdd "-24h" }},0
    total_alarms,{{ nowDate }},3
    total_alarms,{{ nowDateAdd "24h" }},0

    """

  Scenario: given export request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_alarms&sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-filter-to-total-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    total_alarms,{{ parseTime "06-09-2020 00:00" }},0
    total_alarms,{{ parseTime "07-09-2020 00:00" }},0
    total_alarms,{{ parseTime "08-09-2020 00:00" }},0

    """

  Scenario: given export request with filter by entity infos should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_alarms&sampling=day&from={{ nowDateAdd "-3d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-total-alarm-metrics-get-by-entity-infos
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    total_alarms,{{ nowDateAdd "-72h" }},0
    total_alarms,{{ nowDateAdd "-48h" }},0
    total_alarms,{{ nowDateAdd "-24h" }},0
    total_alarms,{{ nowDate }},2
    total_alarms,{{ nowDateAdd "24h" }},0

    """

  Scenario: given export request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "parameters[]": "Parameters is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day&parameters[]=total_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "filter \"not-exist\" not found"
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "parameter \"not-exist\" is not supported"
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_user_activity&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "parameter \"total_user_activity\" is not supported"
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?sampling=not-exist&from={{ now }}&to={{ now }}&parameters[]=total_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  Scenario: given export request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 401

  Scenario: given export request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 403

  Scenario: given export request with all parameters should return all metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_alarms&parameters[]=non_displayed_alarms&parameters[]=instruction_alarms&parameters[]=pbehavior_alarms&parameters[]=correlation_alarms&parameters[]=ack_alarms&parameters[]=cancel_ack_alarms&parameters[]=ack_without_cancel_alarms&parameters[]=ticket_alarms&parameters[]=without_ticket_alarms&parameters[]=ratio_correlation&parameters[]=ratio_instructions&parameters[]=ratio_tickets&parameters[]=ratio_non_displayed&parameters[]=average_ack&parameters[]=average_resolve&sampling=day&from={{ nowDateAdd "-1d" }}&to={{ nowDateAdd "1d" }}&filter=test-filter-to-all-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    total_alarms,{{ nowDateAdd "-24h" }},0
    total_alarms,{{ nowDate }},3
    total_alarms,{{ nowDateAdd "24h" }},0
    non_displayed_alarms,{{ nowDateAdd "-24h" }},0
    non_displayed_alarms,{{ nowDate }},1
    non_displayed_alarms,{{ nowDateAdd "24h" }},0
    instruction_alarms,{{ nowDateAdd "-24h" }},0
    instruction_alarms,{{ nowDate }},1
    instruction_alarms,{{ nowDateAdd "24h" }},0
    pbehavior_alarms,{{ nowDateAdd "-24h" }},0
    pbehavior_alarms,{{ nowDate }},1
    pbehavior_alarms,{{ nowDateAdd "24h" }},0
    correlation_alarms,{{ nowDateAdd "-24h" }},0
    correlation_alarms,{{ nowDate }},1
    correlation_alarms,{{ nowDateAdd "24h" }},0
    ack_alarms,{{ nowDateAdd "-24h" }},0
    ack_alarms,{{ nowDate }},2
    ack_alarms,{{ nowDateAdd "24h" }},0
    cancel_ack_alarms,{{ nowDateAdd "-24h" }},0
    cancel_ack_alarms,{{ nowDate }},1
    cancel_ack_alarms,{{ nowDateAdd "24h" }},0
    ack_without_cancel_alarms,{{ nowDateAdd "-24h" }},0
    ack_without_cancel_alarms,{{ nowDate }},1
    ack_without_cancel_alarms,{{ nowDateAdd "24h" }},0
    ticket_alarms,{{ nowDateAdd "-24h" }},0
    ticket_alarms,{{ nowDate }},1
    ticket_alarms,{{ nowDateAdd "24h" }},0
    without_ticket_alarms,{{ nowDateAdd "-24h" }},0
    without_ticket_alarms,{{ nowDate }},2
    without_ticket_alarms,{{ nowDateAdd "24h" }},0
    ratio_correlation,{{ nowDateAdd "-24h" }},0
    ratio_correlation,{{ nowDate }},33.33
    ratio_correlation,{{ nowDateAdd "24h" }},0
    ratio_instructions,{{ nowDateAdd "-24h" }},0
    ratio_instructions,{{ nowDate }},33.33
    ratio_instructions,{{ nowDateAdd "24h" }},0
    ratio_tickets,{{ nowDateAdd "-24h" }},0
    ratio_tickets,{{ nowDate }},33.33
    ratio_tickets,{{ nowDateAdd "24h" }},0
    ratio_non_displayed,{{ nowDateAdd "-24h" }},0
    ratio_non_displayed,{{ nowDate }},33.33
    ratio_non_displayed,{{ nowDateAdd "24h" }},0
    average_ack,{{ nowDateAdd "-24h" }},0
    average_ack,{{ nowDate }},500
    average_ack,{{ nowDateAdd "24h" }},0
    average_resolve,{{ nowDateAdd "-24h" }},0
    average_resolve,{{ nowDate }},1000
    average_resolve,{{ nowDateAdd "24h" }},0

    """
