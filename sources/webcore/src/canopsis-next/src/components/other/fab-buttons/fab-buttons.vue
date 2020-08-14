<template lang="pug">
  .fab
    v-layout(column)
      refresh-btn(@click="$emit('refresh')")
      v-tooltip(left)
        v-btn(
          slot="activator",
          color="primary",
          fab,
          @click.stop="$emit('create')"
        )
          v-icon add
        slot
</template>

<script>
import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';

export default {
  components: {
    RefreshBtn,
  },
  data() {
    return {
      fab: false,
    };
  },
};
</script>
