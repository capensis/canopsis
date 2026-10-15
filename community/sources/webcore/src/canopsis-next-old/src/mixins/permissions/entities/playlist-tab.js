import { authMixin } from '@/mixins/auth';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

export const permissionsEntitiesPlaylistTabMixin = {
  mixins: [authMixin, entitiesViewGroupMixin],
  methods: {
    getAvailableTabsByIds(tabsIds) {
      return tabsIds.map(id => this.getViewTabById(id));
    },
  },
};
