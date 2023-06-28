import { createNamespacedHelpers } from 'vuex';

import { userToForm, formToUserRequest } from '@/helpers/entities/user/form';

const { mapActions: authMapActions, mapGetters: authMapGetters } = createNamespacedHelpers('auth');
const { mapActions: userMapActions } = createNamespacedHelpers('user');

export const tourBaseMixin = {
  computed: {
    ...authMapGetters(['currentUser']),
  },
  methods: {
    ...authMapActions(['fetchCurrentUser']),
    ...userMapActions({
      updateCurrentUser: 'updateCurrentUser',
    }),

    async finishTourByName(tourName) {
      const userForm = userToForm(this.currentUser);
      const user = formToUserRequest({
        ...userForm,
        ui_tours: {
          ...this.currentUser.ui_tours,

          [tourName]: true,
        },
      });

      await this.updateCurrentUser({ data: user });

      this.fetchCurrentUser();
    },
  },
};
