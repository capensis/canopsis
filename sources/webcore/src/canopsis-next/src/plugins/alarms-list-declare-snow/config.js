import CreateDeclareSnowEvent from './components/modals/create-declare-snow-event.vue';

const MODALS = {
  createDeclareSnowEvent: 'create-declare-snow-event',
};

const WIDGETS_ACTIONS_TYPES = {
  alarmsList: {
    declareSnow: 'declareSnow',
  },
};

const alarmListActionPanelMixin = {
  data() {
    return {
      actionsMap: {
        declareSnow: {
          type: WIDGETS_ACTIONS_TYPES.alarmsList.declareSnow,
          icon: 'report_problem',
          title: this.$t('alarmList.actions.titles.declareSnow'),
          method: this.showDeclareSnowModal,
        },
      },
    };
  },
  methods: {
    showDeclareSnowModal() {
      this.showModal({
        name: 'create-declare-snow-event',
      });
    },
  },
};

export default {
  constants: {
    MODALS,
    WIDGETS_ACTIONS_TYPES,
  },

  alarmListActionPanel: {
    mixins: [alarmListActionPanelMixin],
    computed: {
      actions(actions) {
        const { filteredActionsMap } = this;

        if (filteredActionsMap.declareSnow) {
          const actionsInline = [filteredActionsMap.declareSnow, ...actions.inline];
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

  components: {
    modals: {
      CreateDeclareSnowEvent,
    },
  },
};
