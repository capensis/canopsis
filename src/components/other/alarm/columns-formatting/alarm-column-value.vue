<template lang="pug">
  div
    div(v-if="component", :is="component", :alarm="alarm") {{ component.value }}
    ellipsis(v-else, :text="alarm | get(column.value, columnFilter, '')")
    info-popup-button(v-if="columnName", :column="column.value", :alarm="alarm", :widget="widget")
</template>

<script>
import State from '@/components/other/alarm/columns-formatting/alarm-column-value-state.vue';
import ExtraDetails from '@/components/other/alarm/columns-formatting/alarm-column-value-extra-details.vue';
import InfoPopupButton from '@/components/other/info-popup/popup-button.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';

/**
 * Component to format alarms list columns
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 * @prop {Object} widget - Object representing the widget
 * @prop {Object} column - Property concerned on the column
 */
export default {
  components: {
    State,
    ExtraDetails,
    InfoPopupButton,
    Ellipsis,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
  },
  computed: {
    columnName() {
      return this.column.value.split('.')[1];
    },
    columnFilter() {
      const PROPERTIES_FILTERS_MAP = {
        'v.state.val': value => this.$t(`tables.alarmStatus.${value}`),
        'v.last_update_date': value => this.$d(new Date(value * 1000), 'long'),
        'v.creation_date': value => this.$d(new Date(value * 1000), 'long'),
        'v.last_event_date': value => this.$d(new Date(value * 1000), 'long'),
        'v.state.t': value => this.$d(new Date(value * 1000), 'long'),
        'v.status.t': value => this.$d(new Date(value * 1000), 'long'),
        t: value => this.$d(new Date(value * 1000), 'long'),
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },
    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        'v.state.val': 'state',
        extra_details: 'extra-details',
      };

      return PROPERTIES_COMPONENTS_MAP[this.column.value];
    },
  },
};
</script>
