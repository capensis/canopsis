<template lang="pug">
  v-container
    div
      v-layout(row, wrap)
        v-flex(xs12)
          alarm-list(:widget="widget")
</template>

<script>
import { WIDGET_TYPES } from '@/constants';
import AlarmList from '@/components/other/alarm/alarms-list.vue';
import { generateWidgetByType } from '@/helpers/entities';

export default {
  components: { AlarmList },
  data() {
    const widget = generateWidgetByType(WIDGET_TYPES.alarmList);

    return {
      widget: {
        ...widget,
        parameters: {
          ...widget.parameters,

          widgetColumns: [
            { label: 'Connector name', value: 'v.connector_name' },
          ],
        },
      },
    };
  },
};
</script>
