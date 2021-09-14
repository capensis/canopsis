<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.importExportViews.title') }}
    template(slot="text")
      v-layout(row)
        v-flex.pl-1.pr-1(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.groups') }}
          draggable-groups(
            v-model="importedGroups",
            pull,
            view-pull,
            view-put,
            @change:group="changeImportedGroupHandler"
          )
        v-flex.pl-1(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.views') }}
          draggable-group-views(v-model="importedViews", pull)
        v-divider.ml-1.mr-1.secondary(vertical)
        v-flex(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.result') }}
          draggable-groups(
            v-model="currentGroups",
            put,
            pull,
            view-put,
            @change:group="changeCurrentGroupHandler"
          )
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(
        :loading="submitting",
        type="submit",
        @click="submit"
      ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { get, cloneDeep, isUndefined } from 'lodash';

import { MODALS } from '@/constants';

import {
  viewToRequest,
  groupToRequest,
  groupsWithViewsToPositions,
  prepareImportedViews,
  prepareImportedGroups,
  getExportedGroupsWrappers,
  getExportedViewsWrappersFromGroups,
} from '@/helpers/forms/view';
import { setSeveralFields } from '@/helpers/immutable';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

import { modalInnerMixin } from '@/mixins/modal/inner';
import entitiesViewMixin from '@/mixins/entities/view';
import { authMixin } from '@/mixins/auth';
import { submittableMixinCreator } from '@/mixins/submittable';

import DraggableGroupViews from '@/components/layout/navigation/partial/groups-side-bar/draggable-group-views.vue';
import DraggableGroups from '@/components/layout/navigation/partial/groups-side-bar/draggable-groups.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importExportViews,
  components: {
    DraggableGroupViews,
    DraggableGroups,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    authMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      importedViews: [],
      importedGroups: [],
      currentGroups: [],
    };
  },
  watch: {
    groups: {
      immediate: true,
      handler() {
        this.setDefaultValues();
      },
    },
  },
  methods: {
    changeImportedGroupHandler(groupIndex, group) {
      this.importedGroups.splice(groupIndex, 1, group);
    },

    changeCurrentGroupHandler(groupIndex, group) {
      this.currentGroups.splice(groupIndex, 1, group);
    },

    setDefaultValues() {
      const { importedGroups, importedViews } = this.config;

      this.importedViews = prepareImportedViews(importedViews);
      this.importedGroups = prepareImportedGroups(importedGroups);
      this.currentGroups = cloneDeep(this.groups);
    },

    /**
     * Create exported groups and return updated array
     *
     * @param {ViewGroupWithViews[]} [groups = []]
     * @returns {Promise<ViewGroupWithViews[]>}
     */
    async createGroupsAndGetUpdatedGroups(groups = []) {
      const exportedGroupsWrappers = getExportedGroupsWrappers(groups);

      if (!exportedGroupsWrappers.length) {
        return groups;
      }

      const groupsData = exportedGroupsWrappers
        .map(({ group }) => groupToRequest(group));

      const createdGroups = await this.bulkCreateGroupsWithoutStore({ data: groupsData });

      const groupsFields = createdGroups.reduce((acc, createdGroup, index) => {
        const path = get(exportedGroupsWrappers, [index, 'path']);

        if (!isUndefined(path)) {
          acc[path] = group => ({ ...createdGroup, views: group.views });
        }

        return acc;
      }, {});

      return setSeveralFields(groups, groupsFields);
    },

    /**
     * Create exported views and return updated array
     *
     * @param {ViewGroupWithViews[]} [groups = []]
     * @returns {Promise<ViewGroupWithViews[]>}
     */
    async createViewsAndGetUpdatedGroups(groups) {
      const exportedViewsWrappers = getExportedViewsWrappersFromGroups(groups);

      if (!exportedViewsWrappers.length) {
        return groups;
      }

      const viewData = exportedViewsWrappers
        .map(({ view }) => viewToRequest(view));

      const createdViews = await this.bulkCreateViewsWithoutStore({ data: viewData });

      const viewsFields = createdViews.reduce((acc, createdView, index) => {
        const path = get(exportedViewsWrappers, [index, 'path']);

        if (!isUndefined(path)) {
          acc[path] = createdView;
        }

        return acc;
      }, {});

      return setSeveralFields(groups, viewsFields);
    },

    /**
     * Create exported groups and views and also update positions of all views and groups
     *
     * @returns {Promise}
     */
    async createGroupsAndViews() {
      let groups = await this.createGroupsAndGetUpdatedGroups(this.currentGroups);

      groups = await this.createViewsAndGetUpdatedGroups(groups);

      return this.updateViewsPositions({ data: groupsWithViewsToPositions(groups) });
    },

    async submit() {
      try {
        await this.createGroupsAndViews();
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('errors.default') });
      } finally {
        await this.fetchAllGroupsListWithViewsWithCurrentUser();
        this.$modals.hide();
      }
    },
  },
};
</script>
