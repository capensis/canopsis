import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [rightsTechnicalViewMixin],
  computed: {
    availableGroups() {
      return this.groupsOrdered.reduce((acc, group) => {
        const views = group.views
          .filter(view => this.checkReadAccess(view._id))
          .sort((a, b) => (a.position || Infinity) - (b.position || Infinity));

        if (views.length || this.isEditingMode) {
          acc.push({ ...group, views });
        }

        return acc;
      }, []);
    },
  },
};
