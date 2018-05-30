<template lang="pug">
  span
   template(v-if="attributeType === 'textable'") {{ getProp(alarm, pathToProperty) }}
   template(v-if="attributeType === 'state'")
    state(:stateId="getProp(alarm, pathToProperty)")
   template(v-if="attributeType === 'status'")
    status(:statusId="getProp(alarm, pathToProperty)")
</template>

<script>
import getProp from 'lodash/get';
import State from '@/components/AlarmList/alarm-state-column-value.vue';
import Status from '@/components/AlarmList/alarm-status-column-value.vue';

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
      if (this.pathToProperty === 'v.state.val') return 'state';

      if (this.pathToProperty === 'v.status.val') return 'status';

      return 'textable';
    },
  },
  methods: {
    getProp,
  },
};
</script>
