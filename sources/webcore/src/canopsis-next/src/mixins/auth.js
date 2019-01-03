import { createNamespacedHelpers } from 'vuex';

import { USERS_RIGHTS_MASKS } from '@/constants';
import { checkUserAccess } from '@/helpers/right';

const { mapGetters, mapActions } = createNamespacedHelpers('auth');

export default {
  computed: {
    ...mapGetters(['isLoggedIn', 'currentUser']),
    ...mapGetters({
      currentUserPending: 'pending',
    }),

    checkAccess() {
      return (rightId, rightMask = 1) => checkUserAccess(this.currentUser, rightId, rightMask);
    },

    checkCreateAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.create);
    },

    checkReadAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.read);
    },

    checkUpdateAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.update);
    },

    checkDeleteAccess() {
      return rightId => this.checkAccess(rightId, USERS_RIGHTS_MASKS.delete);
    },
  },
  methods: {
    ...mapActions(['login', 'logout', 'fetchCurrentUser']),
  },
};
