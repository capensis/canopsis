<template lang="pug">
  span {{ text }}
</template>

<script>
import getProp from 'lodash/get';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    propertyKey: {
      type: String,
      default: 'v.status.val',
    },
  },
  computed: {
    statusId() {
      return getProp(this.alarm, this.propertyKey);
    },
    text() {
      return this.$t(`tables.alarmStatus.${this.statusId}`);
    },
  },
};
</script>
