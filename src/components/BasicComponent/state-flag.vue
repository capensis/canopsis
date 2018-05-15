<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('entities/alarmConvention');

export default {
  name: 'state-flag',
  props: {
    val: {
      type: Number,
      default: 0,
    },
    isStatus: {
      type: Boolean,
      default: false,
    },
    isCroppedState: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['getStateAlarmConvention', 'getStatusAlarmConvention']),
    color() {
      if (this.isStatus) {
        return this.getStatusAlarmConvention(this.val)('color');
      } else if (this.isCroppedState){
        return 'black';
      }
      return this.getStateAlarmConvention(this.val)('color');
    },
    icon() {
      if (this.isStatus) {
        return this.getStatusAlarmConvention(this.val)('icon');
      } else if (this.isCroppedState) {
        return 'vertical_align_center';
      }
      return this.getStateAlarmConvention(this.val)('icon');
    },
  },
};
</script>

<style scoped>

</style>
