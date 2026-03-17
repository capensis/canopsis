<template>
  <v-list-item
    :to="link.route"
    :class="link.class"
    class="top-bar-menu-link"
    active-class=""
    @click="click"
    @mouseenter="handleMouseEnter"
  >
    <v-list-item-avatar class="ml-2" size="24" rounded="lg">
      <v-icon :class="{ 'text--secondary': !link.color }" :color="link.color" size="24">
        {{ link.icon }}
      </v-icon>
    </v-list-item-avatar>
    <v-list-item-title :class="{ [`${link.color}--text`]: link.color }">
      <v-layout justify-space-between>
        <span>{{ link.title }}</span>
      </v-layout>
    </v-list-item-title>
    <v-list-item-action v-if="link.links">
      <v-icon>arrow_right</v-icon>
    </v-list-item-action>
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
  setup(props, { emit }) {
    const click = () => {
      if (props.link.links?.length) {
        return;
      }

      props.link.handler?.();
      emit('click');
    };

    const handleMouseEnter = event => emit('mouseenter', event);

    return {
      click,
      handleMouseEnter,
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
