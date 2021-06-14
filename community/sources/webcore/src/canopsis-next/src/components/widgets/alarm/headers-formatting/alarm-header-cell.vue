<template lang="pug">
  component(:is="component.is") {{ component.text }}
</template>

<script>
import { ALARM_ENTITY_FIELDS } from '@/constants';

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
  },
  computed: {
    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        [ALARM_ENTITY_FIELDS.priority]: 'alarm-header-priority',
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
