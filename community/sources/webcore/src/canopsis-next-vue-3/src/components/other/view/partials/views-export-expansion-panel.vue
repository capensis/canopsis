<template>
  <v-expansion-panels
    :model-value="openedPanels"
    class="views-export-expansion-panel"
    multiple
    readonly
    hide-actions
    focusable
  >
    <group-panel
      v-for="group in groups"
      :key="group._id"
      :group="group"
      hide-actions
    >
      <template #title>
        <c-layout align="center" fill>
          <c-checkbox
            :model-value="isGroupSelected(group._id)"
            class="views-export-expansion-panel__checkbox mt-0 pt-0"
            color="primary"
            hide-details
            @update:model-value="changeGroup(group, $event)"
          />
          <span>{{ group.title }}</span>
        </c-layout>
      </template>

      <group-view-panel
        v-for="view in group.views"
        :key="view._id"
        :view="view"
      >
        <template #title>
          <c-layout align="center" justify="space-between" fill>
            <c-checkbox
              :model-value="isViewSelected(view._id)"
              class="views-export-expansion-panel__checkbox"
              color="primary"
              hide-details
              @update:model-value="changeView(view._id, $event)"
            >
              <template #label>
                <span class="text-truncate fill-width">
                  {{ view.title }}
                  <span v-show="view.description" class="ml-1">
                    ({{ view.description }})
                  </span>
                </span>
              </template>
            </c-checkbox>
          </c-layout>
        </template>
      </group-view-panel>
    </group-panel>
  </v-expansion-panels>
</template>

<script>
import { difference } from 'lodash';

import { createRangeArray } from '@/helpers/array';

import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';

export default {
  components: { GroupPanel, GroupViewPanel },
  props: {
    groups: {
      type: Array,
      default: () => [],
    },
    selected: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:selected'],

  computed: {
    openedPanels() {
      return createRangeArray(this.groups.length);
    },
  },

  methods: {
    updateSelected(selected) {
      this.$emit('update:selected', selected);
    },

    isGroupSelected(groupId) {
      return this.selected.groups.includes(groupId);
    },

    isViewSelected(viewId) {
      return this.selected.views.includes(viewId);
    },

    getNextSelectedIds(ids, id, selected) {
      if (selected) {
        return ids.includes(id) ? ids : [...ids, id];
      }

      return ids.filter(selectedId => selectedId !== id);
    },

    changeView(viewId, selected) {
      this.updateSelected({
        ...this.selected,
        views: this.getNextSelectedIds(this.selected.views, viewId, selected),
      });
    },

    changeGroup(group, selected) {
      const groupViews = (group.views ?? []).map(({ _id }) => _id);
      const viewsWithoutGroupViews = difference(this.selected.views, groupViews);

      this.updateSelected({
        groups: this.getNextSelectedIds(this.selected.groups, group._id, selected),
        views: viewsWithoutGroupViews.concat(selected ? groupViews : []),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.views-export-expansion-panel {
  :deep(.group-panel__title) {
    display: flex;
    align-items: center;
  }

  :deep(.panel-item-content__title) {
    display: flex;
    align-items: center;
    min-width: 0;
  }

  &__checkbox {
    flex: 0 0 auto;

    :deep(.v-input__control) {
      min-height: 0;
    }
  }
}
</style>
