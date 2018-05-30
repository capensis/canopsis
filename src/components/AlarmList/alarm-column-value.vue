<template>
  <span>
    <template v-if="attributeType === 'textable'">
      {{ getProp(alarm, pathToProperty) }}
    </template>
    <template v-if="attributeType === 'state'">
      <state :stateId="getProp(alarm, pathToProperty)"></state>
    </template>
  </span>
</template>

<script>
import getProp from 'lodash/get';
import State from '@/components/AlarmList/alarm-state-column-value.vue';

export default {
  components: {
    State,
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

      return 'textable';
    },
  },
  methods: {
    getProp,
  },
};
</script>
