import { createNamespacedHelpers } from 'vuex';

import { setField } from '@/helpers/immutable';
import { prepareUserByData } from '@/helpers/entities';

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
      createUser: 'create',
    }),

    async onStop() {
      if (this.callbacks.onStop) {
        await this.callbacks.onStop();
      }

      const user = prepareUserByData({}, this.currentUser);
      const tours = { ...user.tours, [this.tourName]: true };
      const data = setField(user, 'tours', tours);

      await this.createUser({ data });

      this.fetchCurrentUser();
    },
  },
};
