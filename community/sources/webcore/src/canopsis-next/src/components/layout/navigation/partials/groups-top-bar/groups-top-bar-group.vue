<template>
  <v-menu
    class="group-item"
    content-class="group-v-menu-content bg-secondary"
    location="bottom"
    open-on-hover
  >
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :class="{ 'group-item__active': isActiveGroup }"
        class="group-item__dropdown-btn"
        variant="text"
        color="white"
        @click.stop
      >
        <span class="text-none">{{ group.title }}</span>
        <c-icon v-if="group.is_private" size="16">lock</c-icon>
        <v-btn
          v-show="isNavigationEditingMode"
          class="ml-1"
          variant="text"
          size="small"
          icon
          @click.stop="showEditGroupModal"
        >
          <c-icon size="16">edit</c-icon>
        </v-btn>
        <c-icon class="ml-1">arrow_drop_down</c-icon>
      </v-btn>
    </template>
    <v-list class="bg-secondary">
      <groups-top-bar-group-view v-for="view in group.views" :key="view._id" :view="view" />
    </v-list>
  </v-menu>
</template>

<script>
import { layoutNavigationGroupsBarGroupMixin } from '@/mixins/layout/navigation/groups-bar-group';

import GroupsTopBarGroupView from './groups-top-bar-group-view.vue';

export default {
  components: { GroupsTopBarGroupView },
  mixins: [layoutNavigationGroupsBarGroupMixin],
};
</script>

<style lang="scss" scoped>
.group-item {
  flex: 0 0 auto;

  &__dropdown-btn {
    margin: 6px 2px;
    text-transform: none;
  }

  &__active {
    background: rgb(var(--v-theme-primary));
  }
}
</style>
