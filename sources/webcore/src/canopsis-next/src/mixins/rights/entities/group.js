import { sortBy } from 'lodash';

import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
  computed: {
    availableGroups() {
      return this.groupsOrdered.reduce((acc, group) => {
        const views = group.views.filter(view => this.checkReadAccess(view._id));
        const sortedViews = sortBy(views, ['positions']);

        if (views.length || this.isEditingMode) {
          acc.push({ ...group, views: sortedViews });
        }

        return acc;
      }, []);
    },
  },
};
