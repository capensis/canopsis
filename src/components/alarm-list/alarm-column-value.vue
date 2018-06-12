<template lang="pug">
  span
    template(v-if="attributeType === 'state'")
      state(:stateId="alarm | get(pathToProperty, filter)", :showIcon="showIcon")
    template(v-if="attributeType === 'status'")
      status(:statusId="alarm | get(pathToProperty, filter)")
    template(v-if="attributeType === 'textable'") {{ alarm | get(pathToProperty, filter) }}
</template>

<script>
import State from '@/components/alarm-list/state-column-value.vue';
import Status from '@/components/alarm-list/status-column-value.vue';

export default {
  components: {
    State,
    Status,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    pathToProperty: {
      type: String,
      required: true,
    },
    filter: {
      type: Function,
    },
  },
  computed: {
    attributeType() {
      if (this.pathToProperty === 'v.state.val') {
        return 'state';
      }
      if (this.pathToProperty === 'v.status.val') {
        return 'status';
      }
      return 'textable';
    },
    showIcon() {
      return this.$options.filters.get(this.alarm, 'v.state._t') === 'changestate';
    },
  },
  methods: {
  },
};
</script>
