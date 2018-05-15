<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('entities/alarmConvention');

export default {
  name: 'state-chips',
  props: {
    val: {
      type: Number,
      default: 0,
    },
    isStatus: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['getStateAlarmConvention', 'getStatusAlarmConvention']),
    color() {
      if (this.isStatus) {
        return this.getStatusAlarmConvention(this.val)('color');
      }
      return this.getStateAlarmConvention(this.val)('color');
    },
    text() {
      if (this.isStatus) {
        return this.getStatusAlarmConvention(this.val)('text');
      }
      return this.getStateAlarmConvention(this.val)('text');
    },
  },
};
</script>

<style scoped>

</style>
