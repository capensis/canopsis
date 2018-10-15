<template lang="pug">
  v-container
    div
      v-layout(row, wrap)
        v-flex(xs12)
          alarm-list(:widget="widget")
</template>

<script>
import { WIDGET_TYPES, LIVE_REPORTING_INTERVALS } from '@/constants';
import { generateWidgetByType } from '@/helpers/entities';
import AlarmList from '@/components/other/alarm/alarms-list.vue';
import queryMixin from '@/mixins/query';
import authMixin from '@/mixins/auth';

export default {
  components: { AlarmList },
  mixins: [authMixin, queryMixin],
  data() {
    const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
    const filter = this.$route.query.filter ? JSON.parse(this.$route.query.filter) : null;

    return {
      widget: {
        ...widget,
        parameters: {
          ...widget.parameters,
          mainFilter: filter,
          viewFilters: filter ? [filter] : [],
          widgetColumns: [
            { label: 'Connector name', value: 'v.connector_name' },
          ],
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
