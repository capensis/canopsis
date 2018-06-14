<template lang="pug">
  span
    template(v-if="attributeType === 'textable'") {{ getProp(alarm, pathToProperty) }}
    template(v-if="attributeType === 'state'")
      state(:stateId="getProp(alarm, pathToProperty)", :showIcon="getProp(alarm, 'v.state._t') === 'changestate'")
    template(v-if="attributeType === 'status'")
      status(:statusId="getProp(alarm, pathToProperty)")
    template(v-if="attributeType === 'something'")
      v-tooltip(top)
        v-chip(slot="activator", color="purple")
          v-icon(color="white") check
        span Something

</template>

<script>
import getProp from 'lodash/get';
import State from '@/components/alarm-list/alarm-state-column-value.vue';
import Status from '@/components/alarm-list/alarm-status-column-value.vue';

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
  },
  computed: {
    attributeType() {
      if (this.pathToProperty === 'v.state.val') {
        return 'state';
      }

      if (this.pathToProperty === 'v.status.val') {
        return 'status';
      }

      if (this.pathToProperty === 'something') {
        return 'something';
      }

      return 'textable';
    },
  },
  methods: {
    getProp,
  },
};
</script>
