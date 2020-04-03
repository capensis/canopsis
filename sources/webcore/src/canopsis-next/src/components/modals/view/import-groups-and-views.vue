<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.importExportViews.title') }}
    template(slot="text")
      v-layout(row)
        v-flex.pl-1.pr-1(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.groups') }}
          draggable-groups(v-model="importedGroups", pull, viewPull, viewPut)
        v-flex.pl-1(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.views') }}
          draggable-group-views(v-model="importedViews", pull)
        v-divider.ml-1.mr-1.secondary(vertical)
        v-flex(xs4)
          v-flex.text-xs-center.mb-2 {{ $t('modals.importExportViews.result') }}
          draggable-groups(v-model="currentGroups", put, pull, viewPull, viewPut)
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(
        :loading="submitting",
        type="submit",
        @click="submit"
      ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import entitiesViewsRightsMixin from '@/mixins/entities/view/rights';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';
import entitiesViewMixin from '@/mixins/entities/view';
import submittableMixin from '@/mixins/submittable';

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
    submittableMixin(),
    modalInnerMixin,
    entitiesViewMixin,
    entitiesViewsRightsMixin,
    entitiesViewsGroupsMixin,
    rightsEntitiesGroupMixin,
  ],
  data() {
    return {
      importedViews: [],
      importedGroups: [],
      currentGroups: [],
    };
  },
  computed: {
    groupsOrderedViewsIds() {
      return this.groupsOrdered.reduce((ids, { views }) => {
        ids.push(...views.map(({ _id }) => _id));

        return ids;
      }, []);
    },
  },
  watch: {
    groupsOrdered: {
      immediate: true,
      handler() {
        this.setDefaultValues();
      },
    },
  },
  methods: {
    setDefaultValues() {
      const { groups, views } = this.modal.config;

      this.importedViews = cloneDeep(views);
      this.importedGroups = cloneDeep(groups);
      this.currentGroups = cloneDeep(this.groupsOrdered);
    },

    checkViewsInGroup(views, groupId) {
      return views.reduce((promise, view, viewIndex) => {
        const {
          exported, position, _id: viewId, group_id: viewGroupId, ...viewData
        } = view;
        const data = {
          ...viewData,
          group_id: groupId,
          position: viewIndex,
        };

        if (this.groupsOrderedViewsIds.includes(viewId) && (viewGroupId !== groupId || position !== viewIndex)) {
          return promise.then(() => this.updateViewWithoutStore({ id: viewId, data }));
        } else if (exported) {
          return promise
            .then(() => this.createView({ data }))
            .then(({ _id }) => this.createRightByViewId(_id));
        }

        return promise;
      }, Promise.resolve());
    },

    async updateGroups() {
      await this.currentGroups.reduce((promise, group, groupIndex) => {
        const data = { name: group.name, position: groupIndex };

        if (group.exported) {
          return promise
            .then(() => this.createGroup({ data }))
            .then(({ _id: groupId }) => this.checkViewsInGroup(group.views, groupId));
        } else if (group.position !== groupIndex) {
          return promise.then(this.updateGroup({ data, id: group._id }));
        }

        return promise.then(this.checkViewsInGroup(group.views, group._id));
      }, Promise.resolve());
    },

    async submit() {
      await this.updateGroups();
      this.fetchGroupsList();

      this.$modals.hide();
    },
  },
};
</script>
