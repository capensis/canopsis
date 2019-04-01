<template lang="pug" functional>
  div
    a(
    v-for="(linkUrl, index) in props.link",
    :key="`link-${index}`",
    :href="linkUrl",
    target="_blank",
    @click.stop=""
    ) {{ linkUrl }}
</template>

<script>
export default {
  props: {
    link: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
