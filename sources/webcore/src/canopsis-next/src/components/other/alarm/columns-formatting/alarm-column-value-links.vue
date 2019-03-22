<template lang="pug" functional>
  div
    div(
    v-for="(links, linksKey) in props.alarm.links",
    :key="key"
    )
      a(
      v-for="(link, index) in links",
      :key="`${key}-${index}`",
      :href="link",
      target="_blank",
      @click.stop=""
      ) {{ linksKey }}
</template>

<script>
export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
};
</script>
