import { EVENT_ENTITY_TYPES, MODALS } from '@/constants';

import ModalMixin from './modal/modal';
import EventActionsMixin from './event-actions';

export default {
  mixins: [ModalMixin, EventActionsMixin],
  methods: {
    showActionModal(name) {
      return () => this.showModal({
        name,
        config: this.modalConfig,
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
  },
};
