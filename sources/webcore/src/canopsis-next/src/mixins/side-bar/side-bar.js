import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('sideBar');

export default {
  methods: {
    ...mapActions({
      showSideBar: 'show',
      hideSideBar: 'hide',
    }),
  },
};
