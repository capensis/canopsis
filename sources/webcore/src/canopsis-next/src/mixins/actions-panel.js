import modalMixin from '@/mixins/modal';
import eventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, MODALS, USERS_RIGHTS } from '@/constants';

/**
 * @mixin Mixin for the alarms list actions panel, show modal of the action
 */
export default {
  mixins: [modalMixin, eventActionsMixin],
  methods: {
    showActionModal(name) {
      return () => this.showModal({
        name,
        config: this.modalConfig,
      });
    },

    showMoreInfosModal() {
      this.showModal({
        name: MODALS.moreInfos,
        config: {
          ...this.modalConfig,
          template: this.widget.parameters.moreInfoTemplate,
        },
      });
    },

    showAckRemoveModal() {
      this.showModal({
        name: MODALS.createCancelEvent,
        config: {
          ...this.modalConfig,
          title: 'modals.createAckRemove.title',
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },

    actionsAccessFilterHandler({ type }) {
      const right = USERS_RIGHTS.business.alarmList.actions[type];

      if (!right) {
        return true;
      }

      return this.checkAccess(USERS_RIGHTS.business.alarmList.actions[type]);
    },
  },
};
