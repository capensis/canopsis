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
      :selected-tag="selectedTag"
      @activate="$emit('activate', $event)"
      @select:tag="$emit('select:tag', $event)"
      @clear:tag="$emit('clear:tag')"
      @click:state="$emit('click:state', $event)"
    />
  </color-indicator-wrapper>
  <alarm-column-cell
    v-else
    :alarm="alarm"
    :widget="widget"
    :column="column"
    :small="small"
    :selected-tag="selectedTag"
    @activate="$emit('activate', $event)"
    @select:tag="$emit('select:tag', $event)"
    @clear:tag="$emit('clear:tag')"
    @click:state="$emit('click:state', $event)"
  />
</template>

<script>
import { get } from 'lodash';

import { getAlarmWidgetColumnTemplateId } from '@/helpers/entities/alarm/list';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

import AlarmColumnCell from './alarm-column-cell.vue';

export default {
  components: {
    ColorIndicatorWrapper,
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
    selectedTag: {
      type: String,
      default: '',
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    templateId() {
      return getAlarmWidgetColumnTemplateId(this.widget._id, this.column.value);
    },

    templateContext() {
      return {
        value: get(this.alarm, this.column.value, ''),
        alarm: this.alarm,
        entity: this.alarm.entity,
      };
    },
  },
};
</script>
