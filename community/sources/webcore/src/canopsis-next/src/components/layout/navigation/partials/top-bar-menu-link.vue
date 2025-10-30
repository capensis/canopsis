<template>
  <v-list-item
    :to="link.route"
    class="top-bar-menu-link"
    active-class=""
    @click="click"
  >
    <v-list-item-avatar class="ml-2" size="24">
      <v-icon class="text--secondary" size="24">
        {{ link.icon }}
      </v-icon>
    </v-list-item-avatar>
    <v-list-item-title>
      <v-layout justify-space-between>
        <span>{{ link.title }}</span>
      </v-layout>
    </v-list-item-title>
  </v-list-item>
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const click = () => props.link.handler?.();

    return {
      click,
    };
  },
};
</script>

<style lang="scss" scoped>
.top-bar-menu-link ::v-deep a {
  text-decoration: none;
  color: inherit;

  .v-list-item__avatar {
    min-width: unset;
  }
}
</style>
