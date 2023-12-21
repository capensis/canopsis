<template>
  <component
    v-on="component.on"
    v-bind="component.bind"
    :is="component.is"
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

      return {
        is: component || 'span',
        text: this.header.text,
        bind: {
          style: {
            'white-space': 'normal',
          },
        },
      };
    },
  },
};
</script>
