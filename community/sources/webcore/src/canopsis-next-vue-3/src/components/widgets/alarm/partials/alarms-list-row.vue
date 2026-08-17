<template>
  <tr :data-id="alarm._id" class="alarms-list-row">
    <td
      v-for="header in headers"
      :key="header.key"
      :style="getCellStyle(header)"
      :class="{ 'alarms-list-row__cell--row-icons': header.key === rowIconsColumnKey }"
      class="alarms-list-row__cell"
    >
      <c-layout v-if="header.key === rowIconsColumnKey" class="alarms-list-row__icons" align="center">
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

        <alarms-expand-panel-btn
          v-if="expandable"
          :expanded="expanded"
          :alarm="alarm"
          :widget="widget"
          :small="small"
          :search="search"
          class="alarms-list-row__expand-btn"
          @update:expanded="$emit('update:expanded', $event)"
        />
      </c-layout>

      <c-booted-placeholder-loader
        v-else-if="header.value === 'actions'"
        class="d-flex align-center justify-end"
        async-booting-provider="$asyncBootingActionsPanel"
      >
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
      </c-booted-placeholder-loader>

      <alarm-column-value
        v-else
        :alarm="alarm"
        :widget="widget"
        :column="header"
        :selected-tags="selectedTags"
        :small="small"
        :parent-alarm="parentAlarm"
        @select:tag="$emit('select:tag', $event)"
        @remove:tag="$emit('remove:tag', $event)"
      />

      <span
        v-if="resizing && header.key !== rowIconsColumnKey"
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
import AlarmsExpandPanelBtn from '../expand-panel/alarms-expand-panel-btn.vue';

import AlarmsListRowBookmarkIcon from './alarms-list-row-bookmark-icon.vue';
import AlarmsListRowInstructionsIcon from './alarms-list-row-instructions-icon.vue';

export default {
  components: {
    ActionsPanel,
    AlarmColumnValue,
    AlarmsExpandPanelBtn,
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
    rowIconsColumnKey: {
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
    expandable: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
    search: {
      type: String,
      default: '',
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

  emits: ['update:selected', 'update:expanded', 'start:resize', 'select:tag', 'remove:tag'],

  computed: {
    hasAlarmInstruction() {
      const hasAssignedInstructions = !!this.alarm.assigned_instructions?.length;

      if (this.parentAlarm?.children_instructions && hasAssignedInstructions) {
        return true;
      }

      return hasAssignedInstructions || isNumber(this.alarm.instruction_execution_icon);
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

.alarms-list-row__cell--row-icons {
  padding-left: 6px !important;
  padding-right: 6px !important;
}

.alarms-list-row__checkbox {
  width: 24px;
}

.alarms-list-row__expand-btn {
  margin-left: auto;
}
</style>
