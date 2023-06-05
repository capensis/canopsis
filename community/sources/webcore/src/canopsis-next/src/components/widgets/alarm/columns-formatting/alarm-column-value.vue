<template lang="pug">
  c-runtime-template(v-if="column.template", :template="columnContent")
  color-indicator-wrapper(
    v-else-if="column.colorIndicatorEnabled",
    :type="column.colorIndicator",
    :entity="alarm.entity",
    :alarm="alarm"
  )
    alarm-column-cell(
      :alarm="alarm",
      :widget="widget",
      :column="column",
      :small="small",
      :selected-tag="selectedTag",
      @activate="$emit('activate', $event)",
      @select:tag="$emit('select:tag', $event)"
    )
  alarm-column-cell(
    v-else,
    :alarm="alarm",
    :widget="widget",
    :column="column",
    :small="small",
    :selected-tag="selectedTag",
    @activate="$emit('activate', $event)",
    @select:tag="$emit('select:tag', $event)"
  )
</template>

<script>
import { get } from 'lodash';

import { compile } from '@/helpers/handlebars';

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
  asyncComputed: {
    columnContent: {
      lazy: true,

      async get() {
        const context = {
          value: get(this.alarm, this.column.value, ''),
          alarm: this.alarm,
          entity: this.alarm.entity,
        };
        const compiledTemplate = await compile(this.column.template, context);

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
};
</script>
