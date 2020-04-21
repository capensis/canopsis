<template lang="pug">
  div
    v-layout.pa-4
      v-flex.export-views-block.mr-2(xs6)
        v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
          group-panel(
            v-for="group in groupsOrdered",
            :group="group",
            :key="group._id",
            hideActions
          )
            v-expansion-panel.tabs-panel(
              v-for="view in group.views",
              :key="view._id",
              :value="getPanelValueFromArray(view.tabs)",
              readonly,
              hide-actions,
              expand,
              dark,
              focusable
            )
              v-expansion-panel-content.secondary.lighten-2(
                v-for="tab in view.tabs",
                hide-actions,
                :key="tab._id"
              )
                group-view-panel(slot="header", :view="view")
                span.pl-5 {{ tab.title }}
      v-flex.export-views-block.ml-2(xs6)
        v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
          group-panel(
            v-for="group in groupsOrdered",
            :group="group",
            :key="group._id",
            hideActions
          )
            group-view-panel(
              v-for="view in group.views",
              :key="view._id",
              :view="view"
            )
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import FileSelector from '@/components/forms/fields/file-selector.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';

import entitiesViewMixin from '@/mixins/entities/view';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

export default {
  components: {
    GroupViewPanel,
    FileSelector,
    GroupPanel,
    GroupsSideBarGroup,
  },
  mixins: [
    entitiesViewMixin,
    entitiesViewGroupMixin,
  ],
  computed: {
    openedPanels() {
      return this.getPanelValueFromArray(this.groupsOrdered);
    },
  },
  methods: {
    getPanelValueFromArray(values = []) {
      return new Array(values.length).fill(true);
    },
  },
};
</script>

<style lang="scss" scoped>
  .tabs-panel {
    /deep/ .v-expansion-panel__header {
      padding: 0;
      margin: 0;
    }
    /deep/ .v-expansion-panel__body {
      margin: 12px 0;
    }
  }
</style>
