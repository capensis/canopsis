<template lang="pug">
  div
    div(
    v-for="(link, key) in alarm.links",
    :key="key"
    )
      alarm-column-value-link(:link="link")
</template>

<script>
import AlarmColumnValueLink from './alarm-column-value-link.vue';

export default {
  components: { AlarmColumnValueLink },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
};
</script>
