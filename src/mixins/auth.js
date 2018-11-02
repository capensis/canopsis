import get from 'lodash/get';
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

    hasAccess(rightId, rightMask = 1) {
      const checksum = get(this.currentUser.rights, [rightId, 'checksum'], 0);
      return (checksum & rightMask) === rightMask;
    },
  },
};
