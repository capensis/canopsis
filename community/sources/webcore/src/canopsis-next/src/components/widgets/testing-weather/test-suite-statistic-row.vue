<template lang="pug">
  v-layout(row)
    v-flex(xs6) {{ count }}
    v-flex.text-lowercase(xs6) {{ label }}
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    count: {
      type: [Number, String],
      default: 0,
    },
  },
};
</script>
