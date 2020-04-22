<template lang="pug">
  v-layout.pa-4
    v-flex.export-views-block.mr-2(xs6)
      v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.groups') }}
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
            v-expansion-panel-content(hide-actions)
              group-view-panel(slot="header", :view="view")
              draggable-tabs(:tabs="view.tabs", pull="clone")
    v-flex.export-views-block.ml-2(xs6)
      v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.result') }}
      draggable-tabs(v-field="playlist.tabs", put, pull)
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import FileSelector from '@/components/forms/fields/file-selector.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import DraggableTabs from '@/components/other/playlist/partials/draggable-tabs.vue';

import entitiesViewMixin from '@/mixins/entities/view';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

import TabPanelContent from './partials/tab-panel-content.vue';

export default {
  components: {
    DraggableTabs,
    TabPanelContent,
    GroupViewPanel,
    FileSelector,
    GroupPanel,
    GroupsSideBarGroup,
  },
  mixins: [
    entitiesViewMixin,
    entitiesViewGroupMixin,
  ],
  model: {
    prop: 'playlist',
    event: 'input',
  },
  props: {
    playlist: {
      type: Object,
      required: false,
    },
  },
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
  }
  .tab-panel-item {
    display: flex;
    align-items: center;
    height: 48px;
  }
</style>
