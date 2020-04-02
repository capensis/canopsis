<template lang="pug">
  div
    v-layout
      v-flex.export-views-block.mr-0.ma-4(xs6)
        v-checkbox(v-model="isAllSelected", :label="$t('importExportViews.selectAll')")
        v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
          group-panel(
            v-for="group in groupsOrdered",
            :group="group",
            :key="group._id",
            hideActions
          )
            template(slot="title")
              v-checkbox.group-checkbox(
                v-model="selectedGroupsIds",
                :value="group._id",
                primary,
                @change="changeGroupHandler(group._id, $event)"
              )
              span.group-title {{ group.name }}
            group-view-panel(
              v-for="view in group.views",
              :key="view._id",
              :view="view"
            )
              template(slot="title")
                v-layout(align-center, row, justify-space-between)
                  v-checkbox(v-model="selectedViewIds", :value="view._id")
                  span {{ view.title }}
                  span.ml-1 ({{ view.description }})
      v-flex.btn-group(xs2)
        v-layout(column)
          v-btn(@click="exportViews", :disabled="selectedDataIsEmpty") {{ $t('common.export') }}
          v-btn
            file-selector.view-import-btn(
              ref="fileSelector",
              multiple,
              hide-details,
              @change="importViews"
            )
              span(slot="activator") {{ $t('common.import') }}
</template>

<script>
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';

import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import entitiesViewMixin from '@/mixins/entities/view';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';

import { saveJsonFile } from '@/helpers/files';
import { getFileTextContent } from '@/helpers/file-select';
import { prepareGroupsAndViewsToImport } from '@/helpers/groups';
import { MODALS } from '@/constants';

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
  data() {
    return {
      selectedGroupsIds: [],
      selectedViewIds: [],
    };
  },
  computed: {
    openedPanels() {
      return new Array(this.groupsOrdered.length).fill(true);
    },
    selectedDataIsEmpty() {
      return !this.selectedGroupsIds.length && !this.selectedViewIds.length;
    },
    viewIds() {
      return this.groups.reduce((acc, { views }) => {
        acc.push(...views.map(({ _id }) => _id));
        return acc;
      }, []);
    },
    groupIds() {
      return this.groups.map(({ _id }) => _id);
    },
    isAllSelected: {
      get() {
        return this.groupIds.every(id => this.selectedGroupsIds.includes(id))
          && this.viewIds.every(id => this.selectedViewIds.includes(id));
      },
      set(checked) {
        if (checked) {
          this.selectedGroupsIds = [...this.groupIds];
          this.selectedViewIds = [...this.viewIds];
        } else {
          this.resetSelected();
        }
      },
    },
  },
  methods: {
    async importViews([file]) {
      try {
        const content = await getFileTextContent(file);
        const { groups = [], views = [] } = JSON.parse(content);

        this.$modals.show({
          name: MODALS.importExportViews,
          config: {
            groups,
            views,
          },
        });
      } catch (e) {
        this.$popups.error({ text: this.$t('errors.default') });
      }

      this.$refs.fileSelector.clear();
    },
    exportViews() {
      const exportData = prepareGroupsAndViewsToImport({
        groups: this.selectedGroupsIds.map(this.getGroupById),
        views: this.selectedViewIds.map(this.getViewById),
      });

      saveJsonFile(exportData, `canopsis_groups_views-${new Date().toLocaleString()}`);

      this.resetSelected();
    },
    resetSelected() {
      this.selectedGroupsIds = [];
      this.selectedViewIds = [];
    },
    changeGroupHandler(groupId, selectedGroupsIds) {
      const checked = selectedGroupsIds.includes(groupId);
      const { views } = this.getGroupById(groupId);

      const selectedViewIds = views.map(({ _id }) => _id);
      const viewsIdsWithoutGroupViews = this.selectedViewIds.filter(viewId => !this.viewIds.includes(viewId));

      if (checked) {
        this.selectedViewIds = [
          ...viewsIdsWithoutGroupViews,
          ...selectedViewIds,
        ];
      } else {
        this.selectedViewIds = viewsIdsWithoutGroupViews;
      }
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

  .btn-group {
    display: flex;
    align-items: center;
  }

  .view-import-btn {
    display: inline-flex;
  }

  .group-title {
    overflow: auto;
  }
  .export-views-block {
    & /deep/ .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
  }
</style>
