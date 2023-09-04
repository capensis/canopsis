<template lang="pug">
  v-layout(v-on="$listeners", row, justify-space-around, align-center)
    v-icon(size="18", :style="{ color, caretColor: color }") {{ icon }}
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
