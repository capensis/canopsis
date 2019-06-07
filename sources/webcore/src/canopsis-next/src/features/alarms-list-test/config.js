import CreateTestEvent from './components/modals/create-test-event.vue';

const MODALS = {
  createTestEvent: 'create-test-event',
};

const WIDGETS_ACTIONS_TYPES = {
  alarmsList: {
    test: 'test',
  },
};

const alarmListActionPanelMixin = {
  data() {
    return {
      actionsMap: {
        test: {
          type: WIDGETS_ACTIONS_TYPES.alarmsList.test,
          icon: 'help',
          title: this.$t('alarmList.actions.titles.test'),
          method: this.showTestModal,
        },
      },
    };
  },
  methods: {
    showTestModal() {
      this.showModal({
        name: MODALS.createTestEvent,
      });
    },
  },
};

export default {
  constants: {
    MODALS,
    WIDGETS_ACTIONS_TYPES,
  },
  components: {
    modals: {
      components: {
        CreateTestEvent,
      },
    },
    alarmListActionPanel: {
      mixins: [alarmListActionPanelMixin],
      computed: {
        actions(actions) {
          const { filteredActionsMap } = this;

          if (filteredActionsMap.test) {
            const actionsInline = [filteredActionsMap.test, ...actions.inline];
            const actionsDropDown = [...actions.dropDown];

            if (actionsInline.length > 3) {
              actionsDropDown.unshift(actionsInline.pop());
            }

            return {
              inline: actionsInline,
              dropDown: actionsDropDown,
            };
          }

          return actions;
        },
      },
    },
  },
};
