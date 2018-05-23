import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('popup');

export default {
  methods: {
    ...mapActions({
      addPopup: 'add',
      removePopup: 'remove',
    }),

    addSuccessPopup(popup) {
      this.addPopup({
        popup: { type: 'success', ...popup },
      });
    },
    addInfoPopup(popup) {
      this.addPopup({
        popup: { type: 'info', ...popup },
      });
    },
    addWarningPopup(popup) {
      this.addPopup({
        popup: { type: 'warning', ...popup },
      });
    },
    addErrorPopup(popup) {
      this.addPopup({
        popup: { type: 'error', ...popup },
      });
    },
  },
};
