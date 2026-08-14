<template>
  <tr :data-id="alarm._id" class="alarms-list-row">
    <td
      v-for="header in headers"
      :key="header.key"
      :style="getCellStyle(header)"
      class="alarms-list-row__cell"
    >
      <c-layout v-if="header.key === rowActionsColumnKey" class="alarms-list-row__icons" align="center">
        <c-layout class="alarms-list-row__checkbox" align="center" justify="center">
          <c-light-checkbox
            v-if="selectable"
            :model-value="selected"
            :disabled="!selectableAlarm"
            @update:model-value="$emit('update:selected', $event)"
          />
        </c-layout>

        <c-layout v-if="hasAlarmInstruction" align="center">
          <alarms-list-row-instructions-icon :alarm="alarm" />
        </c-layout>

        <c-layout v-if="alarm.bookmark" align="center">
          <alarms-list-row-bookmark-icon />
        </c-layout>
      </c-layout>

      <c-layout v-else-if="header.value === 'actions'" justify="end" align="center">
        <actions-panel
          v-if="!hideActions"
          :item="alarm"
          :widget="widget"
          :parent-alarm="parentAlarm"
          :refresh-alarms-list="refreshAlarmsList"
          :small="small"
          :inline-count="actionsInlineCount"
          :ignore-media-query="actionsIgnoreMediaQuery"
        />
      </c-layout>

      <alarm-column-value
        v-else
        :alarm="alarm"
        :widget="widget"
        :column="header"
        :selected-tags="selectedTags"
        :small="small"
        @select:tag="$emit('select:tag', $event)"
        @remove:tag="$emit('remove:tag', $event)"
      />

      <span
        v-if="resizing && header.key !== rowActionsColumnKey"
        class="table__resize-handler"
        @pointerdown.stop.prevent="$emit('start:resize', header.value)"
        @mousedown.stop.prevent=""
        @click.stop=""
      />
    </td>
  </tr>
</template>

<script>
import { isNumber } from 'lodash';

import ActionsPanel from '../actions/actions-panel.vue';
import AlarmColumnValue from '../columns-formatting/alarm-column-value.vue';

import AlarmsListRowBookmarkIcon from './alarms-list-row-bookmark-icon.vue';
import AlarmsListRowInstructionsIcon from './alarms-list-row-instructions-icon.vue';

export default {
  components: {
    ActionsPanel,
    AlarmColumnValue,
    AlarmsListRowBookmarkIcon,
    AlarmsListRowInstructionsIcon,
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
    headers: {
      type: Array,
      required: true,
    },
    rowActionsColumnKey: {
      type: String,
      required: true,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    selectedTags: {
      type: Array,
      default: () => [],
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    selectableAlarm: {
      type: Boolean,
      default: false,
    },
    selected: {
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
    hideActions: {
      type: Boolean,
      default: false,
    },
    actionsInlineCount: {
      type: Number,
      required: false,
    },
    actionsIgnoreMediaQuery: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
  },

  emits: ['update:selected', 'start:resize', 'select:tag', 'remove:tag'],

  computed: {
    hasAlarmInstruction() {
      return !!this.alarm.assigned_instructions?.length || isNumber(this.alarm.instruction_execution_icon);
    },
  },

  methods: {
    getCellStyle(header) {
      const width = header.width ?? header.minWidth;

      if (!width) {
        return {};
      }

      return {
        width,
        minWidth: header.minWidth ?? width,
      };
    },
  },
};
</script>

<style lang="scss">
.alarms-list-row__icons {
  min-width: 84px;
}

.alarms-list-row__checkbox {
  width: 24px;
}
</style>
