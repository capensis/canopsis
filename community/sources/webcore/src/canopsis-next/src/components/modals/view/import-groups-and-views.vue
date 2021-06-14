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
import { omit, cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { generateViewTabId, generateWidgetId } from '@/helpers/entities';

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
    groupsOrderedViews() {
      return this.groupsOrdered.reduce((acc, { views }) => {
        acc.push(...views);

        return acc;
      }, []);
    },

    viewsIds() {
      return this.groupsOrderedViews.map(({ _id }) => _id);
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
    prepareImportedTabs(tabs) {
      return tabs.map(tab => ({
        ...tab,
        _id: generateViewTabId(),
        widgets: tab.widgets.map(widget => ({
          ...widget,
          _id: generateWidgetId(widget.type),
        })),
      }));
    },

    prepareImportedViews(views = []) {
      return views.map(view => ({
        ...view,
        tabs: this.prepareImportedTabs(view.tabs),
      }));
    },

    prepareImportedGroups(groups = []) {
      return groups.map(group => ({
        ...group,
        views: this.prepareImportedViews(group.views),
      }));
    },

    changeImportedGroupHandler(groupIndex, group) {
      const groups = [...this.importedGroups];

      groups.splice(groupIndex, 1, group);

      this.importedGroups = groups;
    },

    changeCurrentGroupHandler(groupIndex, group) {
      const groups = [...this.currentGroups];

      groups.splice(groupIndex, 1, group);

      this.currentGroups = groups;
    },

    setDefaultValues() {
      const { groups, views } = this.modal.config;

      this.importedViews = cloneDeep(this.prepareImportedViews(views));
      this.importedGroups = cloneDeep(this.prepareImportedGroups(groups));
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

        if (this.viewsIds.includes(viewId) && (viewGroupId !== groupId || position !== viewIndex)) {
          return promise.then(() => this.updateViewWithoutStore({ id: viewId, data }));
        } else if (exported) {
          return promise
            .then(() => this.createView({ data }))
            .then(({ _id }) => this.createRightByViewId(_id));
        }

        return promise;
      }, Promise.resolve());
    },

    updateGroups() {
      return this.currentGroups.reduce((promise, group, groupIndex) => {
        const data = { ...omit(group, ['views']), name: group.name, position: groupIndex };

        if (group.exported) {
          return promise
            .then(() => this.createGroup({ data }))
            .then(({ _id: groupId }) => this.checkViewsInGroup(group.views, groupId));
        } else if (group.position !== groupIndex) {
          return promise.then(() => this.updateGroup({ data, id: group._id }));
        }

        return promise.then(() => this.checkViewsInGroup(group.views, group._id));
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
