import Vue from 'vue';

export default function (statValue, columnValue) {
  const PROPERTIES_FILTERS_MAP = {
    state_rate: value => Vue.options.filters.percentage(value),
    ack_time_sla: value => Vue.options.filters.percentage(value),
    resolve_time_sla: value => Vue.options.filters.percentage(value),
    time_in_state: value => Vue.options.filters.duration({ value }),
    mtbf: value => Vue.options.filters.duration({ value }),
  };

  if (PROPERTIES_FILTERS_MAP[columnValue]) {
    return PROPERTIES_FILTERS_MAP[columnValue](statValue);
  }

  return statValue;
}
