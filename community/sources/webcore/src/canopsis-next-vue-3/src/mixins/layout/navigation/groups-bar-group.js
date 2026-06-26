import { mapState } from 'pinia';

import { useActiveViewStore } from '@/stores/modules/active-view';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import { layoutNavigationEditingModeMixin } from './editing-mode';

export const layoutNavigationGroupsBarGroupMixin = {
  mixins: [
    layoutNavigationEditingModeMixin,
    entitiesViewGroupMixin,
    permissionsTechnicalViewMixin,
  ],
  props: {
    group: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState(useActiveViewStore, {
      activeView: 'item',
    }),

    isActiveGroup() {
      return this.activeView?.group?._id === this.group._id;
    },
  },
  methods: {
    showEditGroupModal() {
      // TODO: Enable after group modals are migrated.
    },
  },
};
