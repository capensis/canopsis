<template lang="pug">
  v-layout(v-on="$listeners", row, justify-space-around, align-center)
    v-icon(:color="color", size="18") {{ icon }}
    span.caption {{ count }}
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'white',
    },
  },
};
</script>
