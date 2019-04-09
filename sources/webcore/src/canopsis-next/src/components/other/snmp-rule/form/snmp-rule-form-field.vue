<template lang="pug" functional>
  v-layout.mt-3(row, wrap)
    v-flex(v-show="props.label", xs12)
      .body-2 {{ props.label }}
    slot
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
