import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('navigation');

export default {
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
