<template lang="pug">
  color-indicator-wrapper(
    :type="column.colorIndicator",
    :entity="alarm.entity",
    :alarm="alarm"
  )
    alarm-column-cell(
      :alarm="alarm",
      :widget="widget",
      :column="column",
      :columns-filters="columnsFilters",
      @activate="$emit('activate', $event)"
    )
</template>

<script>
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
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
