export const appRemediationInstructionExecutionsPopupsMixin = {
  beforeUnmount() {
    this.leaveFromSimpleManualExecutions();
  },

  computed: {
    hasExecuteInstructionAccess() {
      return false;
    },
  },

  methods: {
    runExecutionsPopups() {
      this.joinToSimpleManualExecutions();
      this.showPausedExecutionsPopup();
    },

    showInstructionPopup() {},

    joinToSimpleManualExecutions() {},

    leaveFromSimpleManualExecutions() {},

    showPausedExecutionsPopup() {
      return Promise.resolve();
    },
  },
};

export default appRemediationInstructionExecutionsPopupsMixin;
