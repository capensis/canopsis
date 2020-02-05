<template lang="pug">
  v-menu.group-item(
    content-class="group-v-menu-content secondary",
    close-delay="0",
    open-on-hover,
    offset-y,
    bottom,
    dark
  )
    div.v-btn.v-btn--flat.theme--dark(
      :data-test="`dropDownButton-group-${group._id}`",
      slot="activator"
    )
      span {{ group.name }}
      v-btn(
        data-test="editGroupButton",
        v-show="isEditingMode",
        depressed,
        small,
        icon,
        @click.stop="showEditGroupModal"
      )
        v-icon(small) edit
      v-icon(dark) arrow_drop_down
    v-list(:data-test="`dropDownZone-group-${group._id}`")
      groups-top-bar-group-view(
        v-for="view in group.views",
        :key="view._id",
        :view="view",
        :isEditingMode="isEditingMode"
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

<style lang="scss" scoped>
  .group-v-menu-content {
    & /deep/ .v-list {
      background-color: inherit;

      .v-list__tile__title {
        height: 28px;
        line-height: 28px;
      }

      .edit-view-button, .duplicate-view-button {
        vertical-align: top;
        margin: 0 0 0 8px;
      }
    }
  }

  .group-item /deep/ .v-menu__activator .v-btn {
    text-transform: none;
  }
</style>
