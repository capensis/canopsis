<template lang="pug">
  div
    v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
      group-panel(
        v-for="group in groupsOrdered",
        :group="group",
        :key="group._id",
        hideActions
      )
        v-layout(row, slot="title")
          v-checkbox.group-checkbox(
            v-model="groupIds",
            :value="group._id",
            primary,
            @change="changeGroupHandler(group._id, $event)"
          )
          | {{ group.name }}
        group-view-panel(
          v-for="view in group.views",
          :key="view._id",
          :view="view"
        )
          template(slot="title")
            v-layout(align-center, row)
              v-checkbox.group-checkbox(v-model="viewIds", :value="view._id")
              | {{ view.name }}
    v-layout(row)
      v-btn(@click="exportViews", :disabled="selectedDataIsEmpty") {{ $t('common.export') }}
      file-selector(
        ref="fileSelector",
        multiple,
        hide-details,
        @change="importViews"
      )
        v-btn(slot="activator") {{ $t('common.import') }}
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
      groupIds: [],
      viewIds: [],
    };
  },
  computed: {
    openedPanels() {
      return new Array(this.groupsOrdered.length).fill(true);
    },
    selectedDataIsEmpty() {
      return !this.groupIds.length && !this.viewIds.length;
    },
  },
  methods: {
    async importViews([file]) {
      try {
        const content = await getFileTextContent(file);
        const { groups, views } = JSON.parse(content);

        this.$modals.show({
          name: MODALS.importExportViews,
          config: {
            groups,
            views,
          },
        });
      } catch (e) {
        this.$popups.error({
          text: this.$t('errors.default'),
        });
      }

      this.$refs.fileSelector.clear();
    },
    exportViews() {
      const exportData = prepareGroupsAndViewsToImport({
        groups: this.groupIds.map(this.getGroupById),
        views: this.viewIds.map(this.getViewById),
      });

      saveJsonFile(exportData, 'groups');

      this.groupIds = [];
      this.viewIds = [];
    },
    changeGroupHandler(groupId, groupIds) {
      const checked = groupIds.includes(groupId);
      const { views } = this.getGroupById(groupId);

      const viewIds = views.map(({ _id }) => _id);
      const viewsIdsWithoutGroupViews = this.viewIds.filter(viewId => !viewIds.includes(viewId));

      if (checked) {
        this.viewIds = [
          ...viewsIdsWithoutGroupViews,
          ...viewIds,
        ];
      } else {
        this.viewIds = viewsIdsWithoutGroupViews;
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

  /deep/ .panel-header span {
    overflow: auto;
  }
</style>
