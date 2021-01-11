<template lang="pug">
  tr(:data-test="`tableRow-${alarm._id}`", :class="{ 'not-filtered': isNotFiltered }")
    td.pr-0(v-if="hasRowActions", data-test="rowCheckbox")
      v-layout(row, align-center, justify-space-between)
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
        v-layout(v-if="hasAlarmInstruction", align-center)
          v-tooltip(top)
            v-icon(slot="activator", size="16", color="black") assignment
            span {{ $t('alarmList.instructionInfoPopup') }}
        v-layout(v-if="expandable", :class="{ 'ml-3': !hasAlarmInstruction }", align-center)
          expand-button(
            :class="expandButtonClass",
            :expanded="row.expanded",
            @expand="showExpandPanel"
          )
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
        :parentAlarm="parentAlarm"
      )
</template>

<script>
import { TOURS } from '@/constants';

import { isResolvedAlarm } from '@/helpers/entities';
import { getStepClass } from '@/helpers/tour';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import ExpandButton from '@/components/other/buttons/expand-button.vue';

import widgetExpandPanelAlarm from '@/mixins/widget/expand-panel/alarm/expand-panel';

export default {
  components: {
    ExpandButton,
    ActionsPanel,
    AlarmColumnValue,
  },
  mixins: [widgetExpandPanelAlarm],
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
    expandable: {
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
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
  },
  computed: {
    alarm() {
      return this.row.item;
    },

    hasRowActions() {
      return this.selectable || this.expandable || this.hasAlarmInstruction;
    },

    hasAlarmInstruction() {
      const { assigned_instructions: assignedInstructions = [] } = this.alarm;

      return assignedInstructions.length;
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

    isNotFiltered() {
      return this.parentAlarm
        && this.parentAlarm.filtered
        && !this.parentAlarm.filtered.includes(this.alarm._id);
    },
  },
  methods: {
    async showExpandPanel() {
      if (!this.row.expanded) {
        await this.fetchAlarmItemWithGroupsAndSteps(this.alarm);
      }

      this.row.expanded = !this.row.expanded;
    },
  },
};
</script>

<style lang="scss" scoped>
  .not-filtered {
    opacity: .4;
    transition: opacity .3s linear;

    &:hover {
      opacity: 1;
    }
  }
</style>
