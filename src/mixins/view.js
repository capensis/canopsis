import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the alarms view
 */
export default {
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
    }),
  },
  computed: {
    ...mapGetters({
      activeView: 'activeItem',
      viewPending: 'pending',
    }),
    widgetWrappers() {
      if (!this.activeView) {
        return [];
      }

      return this.activeView.containerwidget.items;
    },
  },
};
