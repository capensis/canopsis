<template lang="pug">
  span {{ stat }}
</template>

<script>

export default {
  props: {
    name: {
      type: String,
      required: true,
    },
    stats: {
      type: Object,
      required: {},
    },
  },
  computed: {
    stat() {
      if (this.stats && this.stats[this.name]) {
        return this.stats[this.name].sum || 0;
      }

      return 0;
    },
  },
};
</script>
