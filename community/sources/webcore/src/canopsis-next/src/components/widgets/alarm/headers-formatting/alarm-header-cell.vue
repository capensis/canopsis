<template>
  <component
    v-bind="component.bind"
    :is="component.is"
    v-on="component.on"
  >
    {{ component.text }}
  </component>
</template>

<script>
import { ALARM_FIELDS } from '@/constants';

import AlarmHeaderPriority from './alarm-header-priority.vue';
import AlarmHeaderTag from './alarm-header-tag.vue';

export default {
  components: {
    AlarmHeaderPriority,
    AlarmHeaderTag,
  },
  props: {
    header: {
      type: Object,
      required: true,
    },
    selectedTag: {
      type: String,
      default: '',
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
      if (this.header.value === ALARM_FIELDS.tags) {
        return {
          is: 'alarm-header-tag',
          text: this.header.text,
          bind: {
            selectedTag: this.selectedTag,
          },
          on: {
            clear: () => this.$emit('clear:tag'),
          },
        };
      }

      const PROPERTIES_COMPONENTS_MAP = {
        [ALARM_FIELDS.impactState]: 'alarm-header-priority',
      };

      const component = PROPERTIES_COMPONENTS_MAP[this.header.value];
      const bind = this.ellipsisHeaders
        ? { class: 'v-data-table-header-span--ellipsis', title: this.header.text }
        : { style: { 'white-space': 'normal' } };

      return {
        is: component || 'span',
        text: this.header.text,
        bind,
      };
    },
  },
};
</script>
