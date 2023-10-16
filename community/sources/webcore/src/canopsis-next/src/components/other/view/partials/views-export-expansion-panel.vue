<template>
  <v-expansion-panel
    :value="openedPanels"
    readonly
    hide-actions
    expand
    dark
    focusable
  >
    <group-panel
      v-for="group in groups"
      :group="group"
      :key="group._id"
      hide-actions
    >
      <template #title>
        <v-checkbox
          class="group-checkbox mt-0 pt-0"
          :input-value="selected.groups"
          :value="group._id"
          color="primary"
          @change="changeGroup(group, $event)"
        /><span class="group-title">{{ group.title }}</span>
      </template>
      <group-view-panel
        v-for="view in group.views"
        :key="view._id"
        :view="view"
      >
        <template #title>
          <v-layout
            align-center
            justify-space-between
          >
            <v-checkbox
              v-field="selected.views"
              :value="view._id"
              color="primary"
            /><span class="text-truncate">{{ view.title }}<span
              class="ml-1"
              v-show="view.description"
            >({{ view.description }})</span></span>
          </v-layout>
        </template>
      </group-view-panel>
    </group-panel>
  </v-expansion-panel>
</template>

<script>
import { difference } from 'lodash';

import { formMixin } from '@/mixins/form';

import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';

export default {
  components: { GroupPanel, GroupViewPanel },
  mixins: [formMixin],
  model: {
    prop: 'selected',
    event: 'input',
  },
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
  computed: {
    openedPanels() {
      return new Array(this.groups.length).fill(true);
    },
  },
  methods: {
    changeGroup(group, selectedGroups) {
      const checked = selectedGroups.includes(group._id);
      const groupViews = group.views.map(({ _id }) => _id);
      const viewsWithoutGroupViews = difference(this.selected.views, groupViews);

      const selected = {
        groups: selectedGroups,
        views: viewsWithoutGroupViews.concat(checked ? groupViews : []),
      };

      this.updateModel(selected);
    },
  },
};
</script>

<style lang="scss" scoped>
.group-title {
  overflow: auto;
}
</style>
