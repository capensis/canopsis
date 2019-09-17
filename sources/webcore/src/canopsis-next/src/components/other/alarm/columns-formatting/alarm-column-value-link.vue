<template lang="pug" functional>
  div
    a(
      :href="props.link.link",
      target="_blank",
      @click.stop=""
    ) {{ props.link.label }}
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
