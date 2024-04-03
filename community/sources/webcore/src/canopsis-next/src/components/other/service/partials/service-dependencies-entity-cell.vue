<template>
  <component
    v-bind="wrapperProps"
    :is="wrapperProps.is"
  >
    <c-alarm-chip
      v-if="isStateColumn"
      :value="item.entity.state"
    />
    <span v-else>{{ item.entity | get(column.value) }}</span>
  </component>
</template>

<script>
import { ENTITY_FIELDS } from '@/constants';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

export default {
  components: { ColorIndicatorWrapper },
  props: {
    item: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
  },
  computed: {
    wrapperProps() {
      return this.column.colorIndicator
        ? { is: 'color-indicator-wrapper', entity: this.item.entity, type: this.column.colorIndicator }
        : { is: 'span' };
    },

    isStateColumn() {
      return this.column.value?.endsWith(ENTITY_FIELDS.state);
    },
  },
};
</script>
