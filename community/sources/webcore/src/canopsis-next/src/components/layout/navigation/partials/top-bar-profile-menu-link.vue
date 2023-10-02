<template lang="pug">
  v-list-tile(:to="link.route", :class="link.class", active-class="", @click="handleClick")
    v-list-tile-avatar
      v-icon(:color="link.color") {{ link.icon }}
    v-list-tile-title.text-uppercase.body-2(:class="{ [`${link.color}--text`]: link.color }") {{ link.title }}
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    handleClick() {
      this.link.handler?.();
    },
  },
};
</script>
