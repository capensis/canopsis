import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('popup');

/**
 * @mixin Add popup logic, for fetching messages (success, info, warning, error...)
 */
export default {
  methods: {
    ...mapActions({
      addPopup: 'add',
      removePopup: 'remove',
    }),

    addSuccessPopup(popup) {
      this.addPopup({ type: 'success', ...popup });
    },
    addInfoPopup(popup) {
      this.addPopup({ type: 'info', ...popup });
    },
    addWarningPopup(popup) {
      this.addPopup({ type: 'warning', ...popup });
    },
    addErrorPopup(popup) {
      this.addPopup({ type: 'error', ...popup });
    },
  },
};
