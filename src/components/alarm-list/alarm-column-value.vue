<template lang="pug">
  div
    div(v-if="component", :is="component", :alarm="alarm") {{ component.value }}
    span(v-else) {{ alarm | get(property.value, property.filter) }}
</template>

<script>
import State from '@/components/alarm-list/alarm-column-value-state.vue';
import Status from '@/components/alarm-list/alarm-column-value-status.vue';
import Icons from '@/components/alarm-list/alarm-column-value-icons.vue';

const PROPERTIES_COMPONENTS_MAP = {
  'v.state.val': 'state',
  'v.status.val': 'status',
  icons: 'icons',
};

export default {
  components: {
    State,
    Status,
    Icons,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    property: {
      type: Object,
      required: true,
    },
  },
  computed: {
    component() {
      return PROPERTIES_COMPONENTS_MAP[this.property.value];
    },
  },
};
</script>
