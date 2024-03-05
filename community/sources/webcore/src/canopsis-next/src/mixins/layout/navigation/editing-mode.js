import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('navigation');

export const layoutNavigationEditingModeMixin = {
  computed: {
    ...mapGetters({
      isNavigationEditingMode: 'isEditingMode',
    }),
  },
  methods: {
    ...mapActions({
      toggleNavigationEditingMode: 'toggleEditingMode',
    }),
  },
};
