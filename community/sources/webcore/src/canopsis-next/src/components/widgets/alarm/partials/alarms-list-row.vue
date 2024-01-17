<template>
  <tr
    class="alarm-list-row"
    v-on="listeners"
    :class="classes"
  >
    <td
      class="alarm-list-row__icons pr-0"
      v-if="hasRowActions"
    >
      <v-layout
        align-center
        justify-space-between
      >
        <v-layout class="alarm-list-row__checkbox">
          <template v-if="selectable">
            <v-simple-checkbox
              class="ma-0"
              v-if="isAlarmSelectable"
              v-field="selected"
              color="primary"
              hide-details
            />
            <v-simple-checkbox
              v-else
              disabled
              hide-details
            />
          </template>
        </v-layout>
        <v-layout
          v-if="hasAlarmInstruction"
          align-center
        >
          <alarms-list-row-instructions-icon :alarm="alarm" />
        </v-layout>
        <v-layout
          v-if="hasBookmark"
          align-center
        >
          <alarms-list-row-bookmark-icon />
        </v-layout>
        <alarms-expand-panel-btn
          v-if="expandable"
          :expanded="expanded"
          :alarm="alarm"
          :widget="widget"
          :small="small"
          :search="search"
          @input="$emit('expand', $event)"
        />
      </v-layout>
    </td>
    <td
      class="alarm-list-row__cell"
      v-for="header in availableHeaders"
      :key="header.value"
    >
      <actions-panel
        v-if="header.value === 'actions'"
        :item="alarm"
        :widget="widget"
        :parent-alarm="parentAlarm"
        :refresh-alarms-list="refreshAlarmsList"
        :small="small"
        :wrap="wrapActions"
      />
      <alarm-column-value
        v-else
        :alarm="alarm"
        :widget="widget"
        :column="header"
        :selected-tag="selectedTag"
        :small="small"
        @activate="activateRow"
        @select:tag="$emit('select:tag', $event)"
        @click:state="$emit('click:state', $event)"
      />
      <span
        class="alarms-list-table__resize-handler"
        v-if="resizing"
        @mousedown.prevent="$emit('start:resize', header.value)"
        @click.stop=""
      />
    </td>
  </tr>
</template>

<script>
import { flow, isNumber } from 'lodash';

import featuresService from '@/services/features';

import { isActionAvailableForAlarm } from '@/helpers/entities/alarm/form';

import { formBaseMixin } from '@/mixins/form';

import ActionsPanel from '../actions/actions-panel.vue';
import AlarmColumnValue from '../columns-formatting/alarm-column-value.vue';
import AlarmsExpandPanelBtn from '../expand-panel/alarms-expand-panel-btn.vue';

import AlarmsListRowInstructionsIcon from './alarms-list-row-instructions-icon.vue';
import AlarmsListRowBookmarkIcon from './alarms-list-row-bookmark-icon.vue';

export default {
  inject: ['$system'],
  components: {
    ActionsPanel,
    AlarmColumnValue,
    AlarmsExpandPanelBtn,
    AlarmsListRowInstructionsIcon,
    AlarmsListRowBookmarkIcon,
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
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    headers: {
      type: Array,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
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
    resizing: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
    wrapActions: {
      type: Boolean,
      default: false,
    },
    showInstructionIcon: {
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
    hasBookmark() {
      return !!this.alarm.bookmark;
    },

    hasRowActions() {
      return this.selectable || this.expandable || this.showInstructionIcon || this.hasBookmark;
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

    availableHeaders() {
      return this.headers.filter(({ value }) => value);
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

    .v-input--selection-controls__input {
      margin: 0;
    }
  }

  &__icons {
    width: 100px;
  }

  &__cell {
    position: relative;
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
