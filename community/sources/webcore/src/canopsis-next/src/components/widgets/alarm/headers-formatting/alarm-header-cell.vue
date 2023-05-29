<template lang="pug">
  component(
    v-on="component.on",
    v-bind="component.bind",
    :is="component.is"
  ) {{ component.text }}
</template>

<script>
import { ALARM_FIELDS } from '@/constants';

import AlarmHeaderPriority from './alarm-header-priority.vue';
import AlarmHeaderTag from './alarm-header-tag.vue';
import AlarmHeaderActions from './alarm-header-actions.vue';

export default {
  components: {
    AlarmHeaderPriority,
    AlarmHeaderTag,
    AlarmHeaderActions,
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

      if (this.header.value === 'actions') {
        return {
          is: 'alarm-header-actions',
          text: this.header.text,
          bind: {
            resizing: this.resizing,
          },
          on: {
            resizing: () => this.$emit('resizing'),
          },
        };
      }

      const PROPERTIES_COMPONENTS_MAP = {
        [ALARM_FIELDS.impactState]: 'alarm-header-priority',
      };

      const component = PROPERTIES_COMPONENTS_MAP[this.header.value];

      return {
        is: component || 'span',
        text: this.header.text,
      };
    },
  },
};
</script>
