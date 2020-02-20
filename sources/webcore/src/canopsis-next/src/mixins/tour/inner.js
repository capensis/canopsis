import { createNamespacedHelpers } from 'vuex';

import { setField } from '@/helpers/immutable';

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
  },
  mounted() {
    this.$tours[this.tourName].start();
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

      const data = setField(this.currentUser, ['tours', this.tourName], true);

      await this.createUser({ data });

      this.fetchCurrentUser();
    },
  },
};
