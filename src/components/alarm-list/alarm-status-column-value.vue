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
      return this.$t(`tables.alarmStatus.${this.statusId}`);
    },
  },
};
</script>
