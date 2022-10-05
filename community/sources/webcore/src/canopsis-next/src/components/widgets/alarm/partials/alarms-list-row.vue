<template lang="pug">
  tr(v-on="listeners", :class="classes")
    td.pr-0(v-if="hasRowActions")
      v-layout(row, align-center, justify-space-between)
        template(v-if="selectable")
          v-checkbox-functional(
            v-if="!isResolvedAlarm",
            v-field="selected",
            hide-details
          )
          v-checkbox-functional(
            v-else,
            disabled,
            hide-details
          )
        v-layout(v-if="hasAlarmInstruction", align-center)
          alarms-list-row-icon(:alarm="alarm")
        v-layout(v-if="expandable", :class="{ 'ml-3': !hasAlarmInstruction }", align-center)
          alarms-expand-panel-btn(
            v-model="row.expanded",
            :alarm="alarm",
            :widget="widget",
            :is-tour-enabled="isTourEnabled"
          )
    td(v-for="column in columns")
      alarm-column-value(
        :alarm="alarm",
        :widget="widget",
        :column="column",
        :columns-filters="columnsFilters",
        :selected-tag="selectedTag",
        @activate="activateRow",
        @select:tag="$emit('select:tag', $event)"
      )
    td
      actions-panel(
        :item="alarm",
        :widget="widget",
        :is-resolved-alarm="isResolvedAlarm",
        :parent-alarm="parentAlarm",
        :refresh-alarms-list="refreshAlarmsList"
      )
</template>

<script>
import featuresService from '@/services/features';

import { isResolvedAlarm } from '@/helpers/entities';

import ActionsPanel from '../actions/actions-panel.vue';
import AlarmColumnValue from '../columns-formatting/alarm-column-value.vue';
import AlarmsExpandPanelBtn from '../expand-panel/alarms-expand-panel-btn.vue';

import AlarmsListRowIcon from './alarms-list-row-icon.vue';

export default {
  inject: ['$system'],
  components: {
    ActionsPanel,
    AlarmColumnValue,
    AlarmsExpandPanelBtn,
    AlarmsListRowIcon,
  },
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
    columns: {
      type: Array,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    selectedTag: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      active: false,
    };
  },
  computed: {
    ...featuresService.get('components.alarmListRow.computed', {}),

    alarm() {
      return this.row.item;
    },

    hasRowActions() {
      return this.selectable || this.expandable || this.hasAlarmInstruction;
    },

    hasAlarmInstruction() {
      const { children_instructions: parentAlarmChildrenInstructions = false } = this.parentAlarm || {};
      const {
        assigned_instructions: assignedInstructions = [],
        is_auto_instruction_running: isAutoInstructionRunning = false,
        is_manual_instruction_waiting_result: isManualInstructionWaitingResult = false,
        is_all_auto_instructions_completed: isAutoInstructionCompleted = false,
      } = this.alarm;

      const hasAssignedInstructions = !!assignedInstructions.length;

      if (parentAlarmChildrenInstructions && hasAssignedInstructions) {
        return true;
      }

      return hasAssignedInstructions
          || isAutoInstructionRunning
          || isAutoInstructionCompleted
          || isManualInstructionWaitingResult;
    },

    isResolvedAlarm() {
      return isResolvedAlarm(this.alarm);
    },

    isNotFiltered() {
      return this.parentAlarm
        && this.parentAlarm.filtered_children
        && !this.parentAlarm.filtered_children.includes(this.alarm._id);
    },

    listeners() {
      if (featuresService.has('components.alarmListRow.computed.listeners')) {
        return featuresService.call('components.alarmListRow.computed.listeners', this, {});
      }

      return {};
    },

    classes() {
      const classes = { 'not-filtered': this.isNotFiltered, 'grey lighten-3': this.active };

      if (featuresService.has('components.alarmListRow.computed.classes')) {
        return featuresService.call('components.alarmListRow.computed.classes', this, classes);
      }

      return classes;
    },
  },
  methods: {
    activateRow(value) {
      this.active = value;
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
