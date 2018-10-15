<template lang="pug">
  v-container
    div
      v-layout(row, wrap)
        v-flex(xs12)
          alarm-list(:widget="widget")
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { WIDGET_TYPES, LIVE_REPORTING_INTERVALS } from '@/constants';
import { generateWidgetByType } from '@/helpers/entities';
import AlarmList from '@/components/other/alarm/alarms-list.vue';
import queryMixin from '@/mixins/query';
import authMixin from '@/mixins/auth';

export default {
  components: { AlarmList },
  mixins: [authMixin, queryMixin],
  data() {
    const { filter, alarmsStateFilter } = this.$route.query;
    const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
    const filterObject = filter ? JSON.parse(filter) : null;

    const widgetParameters = {
      alarmsStateFilter,
      widgetColumns: widget.parameters.widgetColumns.map(column =>
        ({ label: column.label, value: column.value.replace('alarm.', 'v.') })),
    };

    if (!isEmpty(filterObject)) {
      widgetParameters.mainFilter = filterObject;
      widgetParameters.viewFilters = [filterObject];
    }

    return {
      widget: {
        ...widget,
        parameters: {
          ...widget.parameters,
          ...widgetParameters,
        },
      },
    };
  },

  created() {
    const { tstart, tstop } = this.$route.query;

    const query = {
      tstart,
      tstop,
    };

    if (tstart || tstop) {
      query.interval = LIVE_REPORTING_INTERVALS.custom;
    }

    this.mergeQuery({
      id: this.widget._id,
      query,
    });
  },
};
</script>
