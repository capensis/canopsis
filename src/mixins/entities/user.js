import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('user');

export default {
  methods: {
    ...mapActions({
      editUser: 'edit',
    }),
  },
};
