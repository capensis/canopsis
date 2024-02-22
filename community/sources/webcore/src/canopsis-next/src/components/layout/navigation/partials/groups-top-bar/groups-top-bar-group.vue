<template>
  <v-menu
    class="group-item"
    content-class="group-v-menu-content secondary"
    close-delay="0"
    open-on-hover
    offset-y
    bottom
    dark
  >
    <template #activator="{ on }">
      <v-btn
        class="group-item__dropdown-btn"
        text
        dark
        v-on="on"
      >
        <span>{{ group.title }}</span>
        <v-icon
          v-if="group.is_private"
          color="white"
          small
        >
          lock
        </v-icon>
        <v-btn
          v-show="isNavigationEditingMode"
          depressed
          small
          icon
          @click.stop="showEditGroupModal"
        >
          <v-icon small>
            edit
          </v-icon>
        </v-btn>
        <v-icon
          class="ml-0"
          right
          dark
        >
          arrow_drop_down
        </v-icon>
      </v-btn>
    </template>
    <v-list>
      <groups-top-bar-group-view
        v-for="view in group.views"
        :key="view._id"
        :view="view"
      />
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

<style lang="scss">
.group-item {
  &__dropdown-btn {
    margin: 6px;
    text-transform: none;
  }
}
</style>
