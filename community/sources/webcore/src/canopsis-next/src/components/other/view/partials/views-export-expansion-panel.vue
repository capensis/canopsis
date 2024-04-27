<template>
  <v-expansion-panels
    :value="openedPanels"
    accordion
    readonly
    hide-actions
    multiple
    dark
    focusable
  >
    <group-panel
      v-for="group in groups"
      :key="group._id"
      :group="group"
      hide-actions
    >
      <template #title>
        <v-layout align-center>
          <v-checkbox
            :input-value="selected.groups"
            :value="group._id"
            class="group-checkbox mt-0 pt-0"
            color="primary"
            @change="changeGroup(group, $event)"
          />
          <span>{{ group.title }}</span>
        </v-layout>
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
            >
              <template #label>
                <span class="text-truncate fill-width">
                  {{ view.title }}
                  <span
                    v-show="view.description"
                    class="ml-1"
                  >
                    ({{ view.description }})
                  </span>
                </span>
              </template>
            </v-checkbox>
          </v-layout>
        </template>
      </group-view-panel>
    </group-panel>
  </v-expansion-panels>
</template>

<script>
import { difference } from 'lodash';

import { createRangeArray } from '@/helpers/array';

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
      return createRangeArray(this.groups.length);
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
