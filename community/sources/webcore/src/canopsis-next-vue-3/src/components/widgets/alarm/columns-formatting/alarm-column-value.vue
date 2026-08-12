<template>
  <c-compiled-template
    v-if="column.template"
    :template-id="templateId"
    :template="column.template"
    :context="templateContext"
    class="alarm-column-value"
  />
  <color-indicator-wrapper
    v-else-if="column.colorIndicatorEnabled"
    :type="column.colorIndicator"
    :entity="alarm.entity"
    :alarm="alarm"
  >
    <alarm-column-cell
      :alarm="alarm"
      :widget="widget"
      :column="column"
      :small="small"
      :selected-tags="selectedTags"
      @activate="$emit('activate', $event)"
      @click:state="$emit('click:state', $event)"
      @select:tag="$emit('select:tag', $event)"
      @remove:tag="$emit('remove:tag', $event)"
    />
  </color-indicator-wrapper>
  <alarm-column-cell
    v-else
    :alarm="alarm"
    :widget="widget"
    :column="column"
    :small="small"
    :selected-tags="selectedTags"
    @activate="$emit('activate', $event)"
    @click:state="$emit('click:state', $event)"
    @select:tag="$emit('select:tag', $event)"
    @remove:tag="$emit('remove:tag', $event)"
  />
</template>

<script>
import { get } from 'lodash';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';
import {
  getAlarmWidgetColumnTemplateId,
  getAlarmWidgetGroupColumnTemplateId,
} from '@/helpers/entities/alarm/list';

import AlarmColumnCell from './alarm-column-cell.vue';

export default {
  components: {
    AlarmColumnCell,
    ColorIndicatorWrapper,
  },

  emits: ['activate', 'click:state', 'select:tag', 'remove:tag'],

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
    selectedTags: {
      type: Array,
      default: () => [],
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    templateId() {
      return this.parentAlarm
        ? getAlarmWidgetGroupColumnTemplateId(this.widget._id, this.column.value)
        : getAlarmWidgetColumnTemplateId(this.widget._id, this.column.value);
    },

    templateContext() {
      return {
        alarm: this.alarm,
        entity: this.alarm.entity,
        value: get(this.alarm, this.column.value, ''),
      };
    },
  },
};
</script>
