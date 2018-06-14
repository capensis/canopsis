<template lang="pug">
  span
    template(v-if="attributeType === 'textable'") {{ getProp(alarm, pathToProperty) }}
    template(v-if="attributeType === 'state'")
      state(:stateId="getProp(alarm, pathToProperty)", :showIcon="getProp(alarm, 'v.state._t') === 'changestate'")
    template(v-if="attributeType === 'status'")
      status(:statusId="getProp(alarm, pathToProperty)")
    template(v-if="attributeType === 'something'")
      div(v-if="alarm.v.ack")
        v-tooltip(top)
          v-chip(slot="activator", color="purple")
            v-icon(color="white") check
          span Ack
      div(v-if="alarm.v.ticket")
        v-tooltip(top)
          v-chip(slot="activator", color="blue")
            v-icon(color="white") local_play
          span Ticket
      div(v-if="alarm.v.canceled")
        v-tooltip(top)
          v-chip(slot="activator", color="blue-grey")
            v-icon(color="white") delete
          span Cancel
      div(v-if="alarm.v.snooze")
        v-tooltip(top)
          v-chip(slot="activator", color="pink")
            v-icon(color="white") alarm
          span Snooze
      div(v-if="alarm.pbehaviors.length")
        v-tooltip(top)
          v-chip(slot="activator", color="light-blue")
            v-icon(color="white") merge_type
          span Pbehavior

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
