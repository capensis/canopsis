<template lang="pug">
  tr(:data-test="`tableRow-${alarm._id}`")
    td.pr-0(data-test="rowCheckbox")
      v-layout(row, align-center)
        template(v-if="selectable")
          v-checkbox-functional(
            v-if="!isResolvedAlarm",
            v-field="selected",
            hide-details
          )
          v-checkbox-functional(
            v-else,
            :value="false",
            disabled,
            hide-details
          )
        v-layout.ml-2(align-center)
          v-btn.ma-0(
            :class="expandButtonClass",
            icon,
            small,
            @click="showExpandPanel"
          )
            v-icon {{ row.expanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
    td(v-for="column in columns")
      alarm-column-value(
        :alarm="alarm",
        :column="column",
        :columnFiltersMap="columnFiltersMap",
        :widget="widget"
      )
    td
      actions-panel(
        :item="alarm",
        :widget="widget",
        :isResolvedAlarm="isResolvedAlarm",
        :isEditingMode="isEditingMode"
      )
</template>

<script>
import { TOURS } from '@/constants';

import { isResolvedAlarm } from '@/helpers/entities';
import { getStepClass } from '@/helpers/tour';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';

import widgetExpandPanelAlarmTimeLine from '@/mixins/widget/expand-panel/alarm/time-line';

export default {
  components: {
    ActionsPanel,
    AlarmColumnValue,
  },
  mixins: [widgetExpandPanelAlarmTimeLine],
  model: {
    prop: 'selected',
    event: 'input',
  },
  props: {
    selected: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    row: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    columnFiltersMap: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarm() {
      return this.row.item;
    },
    isResolvedAlarm() {
      return isResolvedAlarm(this.alarm);
    },
    expandButtonClass() {
      if (this.isTourEnabled) {
        return getStepClass(TOURS.alarmsExpandPanel, 1);
      }

      return '';
    },
  },
  methods: {
    async showExpandPanel() {
      if (!this.row.expanded && !this.widget.parameters.moreInfoTemplate) {
        await this.fetchItemWithSteps();
      }

      this.row.expanded = !this.row.expanded;
    },
  },
};
</script>
