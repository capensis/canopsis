import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
  computed: {
    availableGroups() {
      return this.groups.reduce((acc, group) => {
        const views = group.views.filter(view => this.checkReadAccess(view._id));

        if (this.isEditingMode) {
          acc.push({ ...group, views });
        } else if (views.length) {
          acc.push({ ...group, views });
        }

        return acc;
      }, []);
    },
  },
};
