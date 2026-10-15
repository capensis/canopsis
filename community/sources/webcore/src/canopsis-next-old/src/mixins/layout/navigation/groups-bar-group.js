import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import { layoutNavigationEditingModeMixin } from './editing-mode';

const { mapGetters } = createNamespacedHelpers('activeView');

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
    ...mapGetters({
      activeView: 'item',
    }),

    isActiveGroup() {
      return this.activeView?.group?._id === this.group._id;
    },
  },
  methods: {
    showEditGroupModal() {
      this.$modals.show({
        name: MODALS.createGroup,
        config: {
          title: this.$t('modals.group.edit.title'),
          group: this.group,
          deletable: this.group.is_private || this.hasDeleteAnyViewAccess,
          remove: async () => {
            const removeFunction = this.group.is_private
              ? this.removePrivateGroup
              : this.removeGroup;

            await removeFunction({ id: this.group._id });

            await this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
          action: async (data) => {
            const removeFunction = this.group.is_private
              ? this.updatePrivateGroup
              : this.updateGroup;

            await removeFunction({ id: this.group._id, data });

            return this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
        },
      });
    },
  },
};
