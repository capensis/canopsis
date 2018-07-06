<template lang="pug">
  div
    div(v-if="component", :is="component", :alarm="alarm") {{ component.value }}
    span(v-else) {{ alarm | get(property.value, property.filter) }}
    info-popup-button(v-if="columnName", :columnName="columnName", :alarm="alarm")
</template>

<script>
import State from '@/components/other/alarm/columns-formatting/alarm-column-value-state.vue';
import ExtraDetails from '@/components/other/alarm/columns-formatting/alarm-column-value-extra-details.vue';
import InfoPopupButton from '@/components/other/info-popup/popup-button.vue';

const PROPERTIES_COMPONENTS_MAP = {
  'v.state.val': 'state',
  extra_details: 'extra-details',
};

/**
 * Component to format alarms list columns
 *
 * @module alarm
 *
 * @prop {Object} [alarm] - Object representing the alarm
 * @prop {Object} [property] - Property concerned on the column
 */
export default {
  components: {
    State,
    ExtraDetails,
    InfoPopupButton,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    property: {
      type: Object,
      required: true,
    },
  },
  computed: {
    component() {
      return PROPERTIES_COMPONENTS_MAP[this.property.value];
    },
    columnName() {
      return this.property.value.split('.')[1];
    },
  },
};
</script>
