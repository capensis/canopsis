<template lang="pug">
  v-breadcrumbs.pa-0(:items="breadcrumbs")
    template(#item="{ item }")
      v-btn.ma-0(
        :disabled="item.last",
        :loading="item.last && pending",
        small,
        flat,
        @click="$emit('click', item)"
      )
        span.text-none {{ item.text }}
</template>

<script>
export default {
  props: {
    previousMaps: {
      type: Array,
      default: () => [],
    },
    activeMap: {
      type: Object,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    breadcrumbs() {
      const breadcrumbs = this.previousMaps.map(({ name }, index) => ({
        index,
        text: name,
      }));

      if (this.activeMap || this.pending) {
        breadcrumbs.push({
          text: this.activeMap?.name,
          last: true,
        });
      }

      return breadcrumbs;
    },
  },
};
</script>
