<template lang="pug">
  .fab
    v-layout(column)
      refresh-btn(@click="$emit('refresh')")
      v-tooltip(v-if="hasAccess", left)
        v-btn(
          :dark="dark",
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
  props: {
    hasAccess: {
      type: Boolean,
      default: true,
    },
    dark: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      fab: false,
    };
  },
};
</script>
