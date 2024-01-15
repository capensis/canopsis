<template>
  <v-list-item
    class="top-bar-menu-link"
    :to="link.route"
    active-class=""
    @click="handleClick"
  >
    <v-list-item-title>
      <v-layout justify-space-between>
        <span>{{ link.title }}</span>
      </v-layout>
    </v-list-item-title>
    <v-list-item-avatar
      class="ml-2"
      size="24"
    >
      <v-icon
        class="text--secondary"
        size="24"
      >
        {{ link.icon }}
      </v-icon>
    </v-list-item-avatar>
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
  methods: {
    handleClick() {
      this.link.handler?.();
    },
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
