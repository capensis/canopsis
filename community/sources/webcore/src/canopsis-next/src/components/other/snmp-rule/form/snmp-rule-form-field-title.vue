<template lang="pug">
  v-layout.mt-3(row, wrap)
    v-flex(v-show="label", xs12)
      div.body-2 {{ label }}
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
