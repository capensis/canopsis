<template>
  <component v-bind="component.bind" :is="component.is">
    {{ component.text }}
  </component>
</template>

<script>
import { ALARM_FIELDS } from '@/constants';

import AlarmHeaderPriority from './alarm-header-priority.vue';

export default {
  components: {
    AlarmHeaderPriority,
  },

  props: {
    header: {
      type: Object,
      required: true,
    },
    resizing: {
      type: Boolean,
      default: false,
    },
    resizable: {
      type: Boolean,
      default: false,
    },
    ellipsisHeaders: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        [ALARM_FIELDS.impactState]: 'alarm-header-priority',
      };

      const component = PROPERTIES_COMPONENTS_MAP[this.header.value];
      const bind = this.ellipsisHeaders
        ? { class: 'v-data-table-header__span--ellipsis', title: this.header.text }
        : { style: { whiteSpace: 'normal' } };

      return {
        bind,
        is: component || 'span',
        text: this.header.text,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.v-data-table-header__span--ellipsis {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
