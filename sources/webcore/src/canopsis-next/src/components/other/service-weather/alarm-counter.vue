<template lang="pug">
  v-layout.mx-1(row, justify-space-around, align-center)
    v-icon(color="white", size="16") {{ icon }}
    span {{ count }}
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
      required: true,
    },
  },
};
</script>
