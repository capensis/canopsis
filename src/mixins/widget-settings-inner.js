import { createNamespacedHelpers } from 'vuex';

const { mapActions: viewMapActions } = createNamespacedHelpers('view');
const {
  mapActions: userPreferenceMapActions,
  mapGetters: userPreferenceMapGetters,
} = createNamespacedHelpers('userPreference');

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...userPreferenceMapGetters({
      getUserPreferenceByWidget: 'getItemByWidget',
    }),

    userPreference() {
      return this.getUserPreferenceByWidget(this.widget);
    },
  },
  methods: {
    ...viewMapActions({
      createWidget: 'createWidget',
      updateWidget: 'updateWidget',
    }),

    ...userPreferenceMapActions({
      createUserPreference: 'create',
    }),

    closeSettings() {
      this.$emit('closeSettings');
    },
  },
};
