<template lang="pug">
  v-layout.py-4(row)
    v-flex.manage-playlist-tabs.mr-2(xs12)
      v-flex.text-xs-center.mb-2 {{ $t('modals.createPlaylist.groups') }}
      v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
        group-panel(
          v-for="group in groups",
          :group="group",
          :key="group._id",
          hideActions
        )
          template(slot="title")
            v-checkbox.group-checkbox(
              :input-value="selectedGroupsIds",
              :value="group._id",
              @change="selectGroupHandler(group, $event)",
              primary
            )
            span.group-title {{ group.name }}
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
                template(slot="title")
                  v-layout(align-center, row, justify-space-between)
                    v-checkbox(
                      :input-value="selectedViewsIds",
                      :value="view._id",
                      @change="selectViewHandler(view, $event)"
                    )
                    span {{ view.title }}
                    span.ml-1(v-show="view.description") ({{ view.description }})
              tab-panel-content(v-for="tab in view.tabs", :key="tab._id", :tab="tab", hideActions)
                template(slot="title")
                  v-layout.ml-5
                    v-checkbox.tab-checkbox(
                      :input-value="selectedTabsIds",
                      :value="tab._id",
                      @change="selectTabHandler(tab, $event)",
                      primary
                    )
                    span {{ tab.title }}
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';

import TabPanelContent from './partials/tab-panel-content.vue';

export default {
  components: {
    TabPanelContent,
    GroupViewPanel,
    GroupPanel,
    GroupsSideBarGroup,
  },
  model: {
    prop: 'selectedTabs',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
    selectedTabs: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    openedPanels() {
      return this.getPanelValueFromArray(this.groups);
    },

    selectedTabsIds() {
      return this.selectedTabs.map(({ _id }) => _id);
    },

    selectedViewsIds() {
      return this.groups.reduce((acc, { views }) => {
        views.forEach(({ _id: viewId, tabs }) => {
          if (tabs.every(({ _id: tabId }) => this.selectedTabsIds.includes(tabId))) {
            acc.push(viewId);
          }
        });

        return acc;
      }, []);
    },

    selectedGroupsIds() {
      return this.groups.reduce((acc, { _id: groupId, views }) => {
        if (views.every(({ _id: viewId }) => this.selectedViewsIds.includes(viewId))) {
          acc.push(groupId);
        }

        return acc;
      }, []);
    },
  },
  methods: {
    getPanelValueFromArray(values = []) {
      return new Array(values.length).fill(true);
    },

    selectTabHandler(tab, checkedTabs) {
      const checked = checkedTabs.includes(tab._id);

      this.$emit('input', this.getSelectedTabs([tab], checked));
    },

    selectViewHandler(view, checkedViewsIds) {
      const checked = checkedViewsIds.includes(view._id);

      this.$emit('input', this.getSelectedTabs(view.tabs, checked));
    },

    selectGroupHandler(group, selectedGroups) {
      const checked = selectedGroups.includes(group._id);
      const groupTabs = group.views.reduce((acc, view) => {
        acc.push(...view.tabs);
        return acc;
      }, []);

      this.$emit('input', this.getSelectedTabs(groupTabs, checked));
    },

    getSelectedTabs(tabs, checked) {
      const tabsWithoutSelected = this.selectedTabs.filter(({ _id: tabId }) => !tabs.some(({ _id }) => _id === tabId));

      return !checked
        ? tabsWithoutSelected
        : [...tabsWithoutSelected, ...tabs];
    },
  },
};
</script>

<style lang="scss" scoped>
  .manage-playlist-tabs {
    & /deep/ .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
    & /deep/ .v-expansion-panel__body {
      transition: none !important;
    }
  }
  .tabs-panel {
    & /deep/ .v-expansion-panel__header {
      padding: 0;
      margin: 0;
    }
    .tab-checkbox {
      flex: none;
      height: 24px;
      margin: 0;
      padding: 0;
    }
  }
</style>
