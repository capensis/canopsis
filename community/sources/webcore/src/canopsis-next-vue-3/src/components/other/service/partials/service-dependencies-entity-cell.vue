<template>
  <component v-bind="wrapperProps" :is="wrapperProps.is">
    <c-alarm-state-chip v-if="isStateColumn" :value="item.entity.state" small />
    <span v-else>
      {{ value }}
    </span>
  </component>
</template>

<script>
import get from 'lodash/get';

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

    value() {
      return get(this.item.entity, this.column.value, '');
    },
  },
};
</script>
