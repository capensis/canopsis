<template lang="pug">
  v-navigation-drawer(
  :temporary="$mq === 'mobile' || $mq === 'tablet'",
  :value="value",
  stateless,
  clipped,
  right,
  app
  )
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title.white--text.text-xs-center {{ title }}
      v-icon.closeIcon(@click.stop="close", color="white") close
    v-divider
    slot
</template>

<script>
/**
 * Settings wrapper component
 *
 * @prop {bool} value - controls visibility of current component
 * @prop {string} title - title for settings header
 */
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    title: String,
  },
  methods: {
    close() {
      this.$emit('input', false);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
