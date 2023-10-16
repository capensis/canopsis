<template>
  <v-layout>
    <v-flex
      class="manage-playlist-tabs mr-2"
      xs12="xs12"
    >
      <v-flex class="text-center mb-2">
        {{ $t('modals.createPlaylist.groups') }}
      </v-flex>
      <v-expansion-panel
        readonly="readonly"
        hide-actions="hide-actions"
        expand="expand"
        dark="dark"
        focusable="focusable"
        :value="openedPanels"
      >
        <group-panel
          v-for="group in groups"
          :group="group"
          :key="group._id"
          hide-actions="hide-actions"
        >
          <template #title="">
            <v-checkbox
              class="group-checkbox mt-0 pt-0"
              :input-value="selectedGroupsIds"
              :value="group._id"
              :disabled="isDisabledGroup(group)"
              color="primary"
              @change="selectGroupHandler(group, $event)"
            /><span class="group-title">{{ group.title }}</span>
          </template>
          <v-expansion-panel
            class="tabs-panel"
            v-for="view in group.views"
            :key="view._id"
            :value="getPanelValueFromArray(view.tabs)"
            readonly="readonly"
            hide-actions="hide-actions"
            expand="expand"
            dark="dark"
            focusable="focusable"
          >
            <v-expansion-panel-content hide-actions="hide-actions">
              <template #header="">
                <group-view-panel :view="view">
                  <template #title="">
                    <v-layout
                      align-center="align-center"
                      justify-space-between="justify-space-between"
                    >
                      <v-checkbox
                        class="group-checkbox mt-0 pt-0"
                        :input-value="selectedViewsIds"
                        :value="view._id"
                        :disabled="isDisabledView(view)"
                        color="primary"
                        @change="selectViewHandler(view, $event)"
                      /><span class="text-truncate">{{ view.title }}<span
                        class="ml-1"
                        v-show="view.description"
                      >({{ view.description }})</span></span>
                    </v-layout>
                  </template>
                </group-view-panel>
              </template>
              <tab-panel-content
                v-for="tab in view.tabs"
                :key="tab._id"
                :tab="tab"
                hide-actions="hide-actions"
              >
                <template #title="">
                  <v-layout
                    class="ml-5"
                    align-center="align-center"
                  >
                    <v-checkbox
                      class="tab-checkbox group-checkbox"
                      :input-value="selectedTabsIds"
                      :value="tab._id"
                      color="primary"
                      @change="selectTabHandler(tab, $event)"
                    /><span>{{ tab.title }}</span>
                  </v-layout>
                </template>
              </tab-panel-content>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </group-panel>
      </v-expansion-panel>
    </v-flex>
  </v-layout>
</template>

<script>
import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';

import TabPanelContent from '../partials/tab-panel-content.vue';

export default {
  components: {
    TabPanelContent,
    GroupViewPanel,
    GroupPanel,
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
        views.forEach(({ _id: viewId, tabs = [] }) => {
          if (tabs.length && tabs.every(({ _id: tabId }) => this.selectedTabsIds.includes(tabId))) {
            acc.push(viewId);
          }
        });

        return acc;
      }, []);
    },

    selectedGroupsIds() {
      return this.groups.reduce((acc, { _id: groupId, views }) => {
        if (views.length && views.every(({ _id: viewId }) => this.selectedViewsIds.includes(viewId))) {
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

      this.updateSelectedTabs([tab], checked);
    },

    selectViewHandler(view, checkedViewsIds) {
      const checked = checkedViewsIds.includes(view._id);

      this.updateSelectedTabs(view.tabs, checked);
    },

    selectGroupHandler(group, selectedGroups) {
      const checked = selectedGroups.includes(group._id);
      const groupTabs = group.views.reduce((acc, view) => {
        acc.push(...view.tabs);
        return acc;
      }, []);

      this.updateSelectedTabs(groupTabs, checked);
    },

    updateSelectedTabs(tabs, checked) {
      const tabIds = tabs.map(({ _id }) => _id);
      const tabsWithoutSelected = this.selectedTabs.filter(({ _id: tabId }) => !tabIds.includes(tabId));

      this.$emit('input', !checked ? tabsWithoutSelected : [...tabsWithoutSelected, ...tabs]);
    },

    isDisabledGroup(group) {
      return !group.views.length || group.views.every(this.isDisabledView);
    },

    isDisabledView(view) {
      return !view.tabs?.length;
    },
  },
};
</script>

<style lang="scss" scoped>
  .manage-playlist-tabs {
    & ::v-deep .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
    & ::v-deep .v-expansion-panel__body {
      transition: none !important;
    }
  }
  .tabs-panel {
    & ::v-deep .v-expansion-panel__header {
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
