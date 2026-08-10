import { mapState } from 'pinia';

import { MODALS, ROUTES_NAMES } from '@/constants';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { useActiveViewStore } from '@/stores/modules/active-view';

import { layoutNavigationEditingModeMixin } from './editing-mode';

export const layoutNavigationGroupsBarGroupMixin = {
  mixins: [layoutNavigationEditingModeMixin, entitiesViewGroupMixin, permissionsTechnicalViewMixin],
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
      return (
        this.$route.name === ROUTES_NAMES.view && this.activeView?.group?._id === this.group._id
      );
    },
  },
  methods: {
    showEditGroupModal() {
      const [removeFunction, updateFunction] = this.group.is_private
        ? [this.removePrivateGroup, this.updatePrivateGroup]
        : [this.removeGroup, this.updateGroup];

      this.$modals.show({
        name: MODALS.createGroup,
        config: {
          title: this.$t('modals.group.edit.title'),
          group: this.group,
          deletable: this.group.is_private || this.hasDeleteAnyViewAccess,
          remove: async () => {
            await removeFunction({ id: this.group._id });
            await this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
          action: async (data) => {
            await updateFunction({ id: this.group._id, data });
            return this.fetchAllGroupsListWithWidgetsWithCurrentUser();
          },
        },
      });
    },
  },
};
