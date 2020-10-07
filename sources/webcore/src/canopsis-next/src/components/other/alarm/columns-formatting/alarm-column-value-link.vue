<template lang="pug">
  div
    a(
      :href="link.link",
      target="_blank",
      @click.stop=""
    ) {{ link.label }}
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
