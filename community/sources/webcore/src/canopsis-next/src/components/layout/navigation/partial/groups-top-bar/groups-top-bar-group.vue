<template lang="pug">
  v-menu.group-item(
    content-class="group-v-menu-content secondary",
    close-delay="0",
    open-on-hover,
    offset-y,
    bottom,
    dark
  )
    template(#activator="{ on }")
      v-btn.group-item__dropdown-btn(v-on="on", flat, dark)
        span {{ group.title }}
        v-btn(
          v-show="isNavigationEditingMode",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal"
        )
          v-icon(small) edit
        v-icon.ml-0(right, dark) arrow_drop_down
    v-list
      groups-top-bar-group-view(
        v-for="view in group.views",
        :key="view._id",
        :view="view"
      )
</template>

<script>
import layoutNavigationGroupsBarGroupMixin from '@/mixins/layout/navigation/groups-bar-group';

import GroupsTopBarGroupView from './groups-top-bar-group-view.vue';

export default {
  components: { GroupsTopBarGroupView },
  mixins: [layoutNavigationGroupsBarGroupMixin],
};
</script>

<style lang="scss">
.group-item {
  &__dropdown-btn {
    text-transform: none;
  }
}
</style>
