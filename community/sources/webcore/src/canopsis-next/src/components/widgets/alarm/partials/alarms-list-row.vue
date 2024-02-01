<template lang="pug">
  tr.alarm-list-row(v-on="listeners", :class="classes")
    td.pr-0(v-if="hasRowActions")
      v-layout(row, align-center, justify-space-between)
        v-layout.alarm-list-row__checkbox
          template(v-if="selectable")
            v-checkbox-functional.ma-0(v-if="isAlarmSelectable", v-field="selected", hide-details)
            v-checkbox-functional(v-else, disabled, hide-details)
        v-layout(v-if="hasAlarmInstruction", align-center, justify-center)
          alarms-list-row-icon(:alarm="alarm")
        alarms-expand-panel-btn(
          v-if="expandable",
          v-model="row.expanded",
          :alarm="alarm",
          :widget="widget",
          :is-tour-enabled="isTourEnabled",
          :small="small",
          :search="search"
        )
    td(v-for="column in columns")
      alarm-column-value(
        :alarm="alarm",
        :widget="widget",
        :column="column",
        :selected-tag="selectedTag",
        :small="small",
        @activate="activateRow",
        @select:tag="$emit('select:tag', $event)"
      )
    td(v-if="!hideActions")
      actions-panel(
        :item="alarm",
        :widget="widget",
        :parent-alarm="parentAlarm",
        :refresh-alarms-list="refreshAlarmsList",
        :small="small"
      )
</template>

<script>
import { flow, isNumber } from 'lodash';

import featuresService from '@/services/features';

import { isActionAvailableForAlarm } from '@/helpers/entities';

import { formBaseMixin } from '@/mixins/form';

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
  mixins: [formBaseMixin],
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
    selecting: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    selectedTag: {
      type: String,
      default: '',
    },
    medium: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    search: {
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
    alarm() {
      return this.row.item;
    },

    hasRowActions() {
      return this.selectable || this.expandable || this.hasAlarmInstruction;
    },

    hasAlarmInstruction() {
      const { children_instructions: parentAlarmChildrenInstructions = false } = this.parentAlarm || {};
      const { assigned_instructions: assignedInstructions = [] } = this.alarm;

      const hasAssignedInstructions = !!assignedInstructions.length;

      if (parentAlarmChildrenInstructions && hasAssignedInstructions) {
        return true;
      }

      return hasAssignedInstructions || isNumber(this.alarm.instruction_execution_icon);
    },

    isAlarmSelectable() {
      return isActionAvailableForAlarm(this.alarm);
    },

    isNotFiltered() {
      return this.alarm.filtered === false;
    },

    listeners() {
      let listeners = {};

      if (featuresService.has('components.alarmListRow.computed.listeners')) {
        listeners = featuresService.call('components.alarmListRow.computed.listeners', this, {});
      }

      if (this.selecting) {
        listeners.mousedown = flow([this.mouseSelecting, listeners.mouseenter].filter(Boolean));
      }

      return listeners;
    },

    classes() {
      const classes = { 'alarm-list-row--not-filtered': this.isNotFiltered, 'grey lighten-3': this.active };

      if (featuresService.has('components.alarmListRow.computed.classes')) {
        return featuresService.call('components.alarmListRow.computed.classes', this, classes);
      }

      return classes;
    },
  },
  methods: {
    mouseSelecting(event) {
      if (event.ctrlKey && event.buttons) {
        event.preventDefault();

        this.updateModel(!this.selected);
      }

      return event;
    },

    activateRow(value) {
      this.active = value;
    },
  },
};
</script>

<style lang="scss">
.alarm-list-row {
  &__checkbox {
    width: 24px;
    max-width: 24px;
    height: 24px;
  }

  &--not-filtered {
    opacity: .4;
    transition: opacity .3s linear;

    &:hover {
      opacity: 1;
    }
  }
}
</style>
