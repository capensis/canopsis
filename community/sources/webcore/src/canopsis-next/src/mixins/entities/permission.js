import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('permission');

export const entitiesPermissionsMixin = {
  methods: {
    ...mapActions({
      fetchPermissionsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
