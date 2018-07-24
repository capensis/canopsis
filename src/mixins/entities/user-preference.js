import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('userPreference');

/**
 * @mixin Helpers for the userPreference entity
 */
export default {
  computed: {
    ...mapGetters({
      getUserPreferenceByWidget: 'getItemByWidget',
    }),

    userPreference() {
      return this.getUserPreferenceByWidget(this.widget);
    },
  },
  methods: {
    ...mapActions({
      fetchUserPreferencesList: 'fetchList',
      fetchUserPreferenceItemByWidgetId: 'fetchItemByWidgetId',
      createUserPreference: 'create',
    }),
  },
};
