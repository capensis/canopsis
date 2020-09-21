import { sortBy } from 'lodash';

import rightsTechnicalViewMixin from '../technical/view';
import layoutNavigationEditingModeMixin from '../../layout/navigation/editing-mode';

export default {
  mixins: [
    rightsTechnicalViewMixin,
    layoutNavigationEditingModeMixin,
  ],
  computed: {
    availableGroups() {
      return this.groupsOrdered.reduce((acc, group) => {
        const views = group.views.filter(view => this.checkReadAccess(view._id));
        const sortedViews = sortBy(views, ['position']);

        if (views.length || this.isNavigationEditingMode) {
          acc.push({ ...group, views: sortedViews });
        }

        return acc;
      }, []);
    },
  },
};
