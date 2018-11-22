<template lang="pug">
  v-layout.secondary.groups-wrapper
    v-toolbar-items
      v-carousel(
      :cycle="false",
      :height="48",
      max="500",
      hide-delimiters
      )
        v-carousel-item(
        v-for="group in groups",
        :key="group._id"
        )
          v-menu(
          content-class="group-v-menu-content secondary",
          close-delay="0",
          bottom,
          open-on-hover,
          offset-y,
          dark
          )
            v-btn(slot="activator", flat, dark) {{ group.name }}
              v-icon(dark) arrow_drop_down
            v-list
              v-list-tile(
              v-for="view in group.views",
              :key="view._id",
              :to="{ name: 'view', params: { id: view._id } }",
              )
                v-list-tile-title {{ view.title }}
</template>

<script>
import EntitiesViewGroupMixin from '@/mixins/entities/view/group';

export default {
  mixins: [EntitiesViewGroupMixin],
  mounted() {
    this.fetchGroupsList();
  },
};
</script>

<style lang="scss">
  .groups-wrapper {
    height: 48px;

    .v-menu__activator .v-btn {
      text-transform: none;
    }
  }

  .group-v-menu-content {
    .v-list {
      background-color: inherit;
    }
  }
</style>
