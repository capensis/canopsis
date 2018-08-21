import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('auth');

export default {
  computed: {
    ...mapGetters(['isLoggedIn', 'currentUser']),
    ...mapGetters({
      currentUserPending: 'pending',
    }),
  },
  methods: {
    ...mapActions(['login', 'logout', 'fetchCurrentUser']),
  },
};
