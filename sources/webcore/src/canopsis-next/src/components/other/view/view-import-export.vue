<template lang="pug">
  div
    v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
      group-panel(
        v-for="group in availableGroups",
        :group="group",
        :key="group._id",
        hideActions
      )
        v-layout(align-center, slot="title")
          v-checkbox-functional(:value="group._id", primary, @change="changeGroupHandler")
          | {{ group.name }}
        group-view-panel(
          v-for="view in group.views",
          :key="view._id",
          :view="view"
        )
          template(slot="title")
            v-layout(align-center, row)
              v-checkbox-functional.group-checkbox(:value="view._id", @change="changeViewHandler")
              | {{ view.name }}
    v-btn(@click="exportViews") {{ $t('common.export') }}
    v-btn(@click="importViews") {{ $t('common.import') }}
</template>

<script>
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';

import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';

export default {
  components: { GroupViewPanel, GroupPanel, GroupsSideBarGroup },
  mixins: [
    rightsEntitiesGroupMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    return {
      selectedGroups: [],
      selectedViews: [],
    };
  },
  computed: {
    openedPanels() {
      return new Array(this.availableGroups.length).fill(true);
    },
  },
  methods: {
    importViews() {},
    exportViews() {},
    changeGroupHandler(groupId) {
      this.selectedGroups.push(groupId);
    },
    changeViewHandler(viewId) {
      this.selectedViews.push(viewId);
    },
  },
};
</script>

<style lang="scss">
  .group-checkbox {
    height: 24px;
    margin: 0;
    padding: 0;
  }

  /deep/ .panel-header span {
    overflow: auto;
  }
</style>
