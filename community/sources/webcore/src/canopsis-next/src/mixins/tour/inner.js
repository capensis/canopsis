import { createNamespacedHelpers } from 'vuex';

import { formToUser, userToForm } from '@/helpers/forms/user';

const { mapActions: authMapActions, mapGetters: authMapGetters } = createNamespacedHelpers('auth');
const { mapActions: userMapActions } = createNamespacedHelpers('user');

export default {
  props: {
    callbacks: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...authMapGetters(['currentUser']),

    tourCallbacks() {
      return {
        ...this.callbacks,

        onStop: this.onStop,
      };
    },
    tourInstance() {
      return this.$tours[this.tourName];
    },
  },
  mounted() {
    if (this.tourInstance) {
      this.tourInstance.start();
    }
  },
  methods: {
    ...authMapActions(['fetchCurrentUser']),
    ...userMapActions({
      updateUser: 'update',
    }),

    async onStop() {
      if (this.callbacks.onStop) {
        await this.callbacks.onStop();
      }

      const userForm = userToForm(this.currentUser);
      const user = formToUser({
        ...userForm,
        ui_tours: {
          ...this.currentUser.ui_tours,
          [this.tourName]: true,
        },
      });

      await this.updateUser({ data: user, id: user._id });

      this.fetchCurrentUser();
    },
  },
};
