import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
  computed: {
    availableGroups() {
      return this.groups.reduce((acc, group) => {
        const views = group.views.filter(view => this.checkReadAccess(view._id));

        if (views.length || this.isEditingMode) {
          acc.push({ ...group, views });
        }

        return acc;
      }, []);
    },
  },
};
