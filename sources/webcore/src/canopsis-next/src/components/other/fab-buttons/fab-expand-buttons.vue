<template lang="pug">
  .fab
    v-layout(column)
      refresh-btn(@click="$emit('refresh')")
      v-speed-dial(
        v-if="hasAccess",
        v-model="fab",
        direction="left",
        :transition="transition"
      )
        v-btn(slot="activator", color="primary", fab, v-model="fab")
          v-icon add
          v-icon close
        slot
</template>

<script>
import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';

export default {
  components: {
    RefreshBtn,
  },
  props: {
    transition: {
      type: String,
      default: 'slide-y-reverse-transition',
    },
    hasAccess: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      fab: false,
    };
  },
};
</script>
