<template lang="pug">
  span {{ text }}
</template>

<script>
export default {
  props: {
    statusId: {
      required: true,
    },
  },
  computed: {
    text() {
      return this.$t(`tables.alarmStatus.${this.stringId}`);
    },
    stringId() {
      // todo use config as for states
      switch (this.statusId) {
        case 0:
          return 'off';
        case 1:
          return 'ongoing';
        case 2:
          return 'flapping';
        case 3:
          return 'stealthy';
        case 4:
          return 'canceled';
        default:
          return 'unknown';
      }
    },
  },
};
</script>

<style scoped>

</style>
