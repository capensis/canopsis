<template lang="pug">
  div(:data-test="`alarmValue-${column.text}`")
    v-tooltip(v-if="column.isState", right)
      div(slot="activator", :style="stateStyle", :class="{ 'state-column-wrapper': column.isState }")
        alarm-column-cell(
          :alarm="alarm",
          :widget="widget",
          :column="column",
          :columnFiltersMap="columnFiltersMap"
        )
      span {{ stateData.text }}
    alarm-column-cell(
      v-else,
      :alarm="alarm",
      :widget="widget",
      :column="column",
      :columnFiltersMap="columnFiltersMap"
    )
</template>

<script>
import { ALARM_ENTITY_FIELDS } from '@/constants';
import { formatState } from '@/helpers/formatting';

import AlarmColumnCell from '@/components/widgets/alarm/columns-formatting/alarm-column-cell.vue';

export default {
  components: {
    AlarmColumnCell,
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
    columnFiltersMap: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isStateField() {
      return this.column.value === ALARM_ENTITY_FIELDS.state;
    },

    stateData() {
      return formatState(this.alarm.v.state.val);
    },

    stateStyle() {
      return this.column.isState && !this.isStateField
        ? { backgroundColor: this.stateData.color }
        : {};
    },
  },
};
</script>

<style lang="scss" scoped>
  .state-column-wrapper {
    display: inline-block;
    border-radius: 10px;
    padding: 3px 7px;
  }
</style>
