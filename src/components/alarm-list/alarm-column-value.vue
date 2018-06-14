<template lang="pug">
  span
    template(v-if="attributeType === 'state'")
      state(:stateId="alarm | get(pathToProperty, filter)", :showIcon="showIcon")
    template(v-if="attributeType === 'status'")
      status(:statusId="alarm | get(pathToProperty, filter)")
    template(v-if="attributeType === 'textable'") {{ alarm | get(pathToProperty, filter) }}
    info-popup-button(:columnName="columnName", :alarm="alarm")
</template>

<script>
import State from '@/components/alarm-list/alarm-state-column-value.vue';
import Status from '@/components/alarm-list/alarm-status-column-value.vue';
import InfoPopupButton from '@/components/other/info-popup/popup-button.vue';

export default {
  components: {
    State,
    Status,
    InfoPopupButton,
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
    columnName() {
      return this.pathToProperty.split('.')[1];
    },
  },
  methods: {
  },
};
</script>
