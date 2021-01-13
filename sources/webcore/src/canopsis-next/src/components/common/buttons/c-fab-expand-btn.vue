<template lang="pug">
  .fab
    v-layout(column)
      c-refresh-btn(@click="$emit('refresh')")
      v-speed-dial(
        v-if="hasAccess",
        v-model="fab",
        direction="left",
        :transition="transition"
      )
        v-btn(slot="activator", color="primary", fab, v-model="fab")
          slot(name="icons")
            v-icon add
            v-icon close
        slot
</template>

<script>
export default {
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
