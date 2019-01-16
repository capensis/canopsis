import omit from 'lodash/omit';

import { MODALS, EVENT_ENTITY_TYPES, BUSINESS_USER_RIGHTS_ACTIONS_MAP } from '@/constants';

import modalMixin from '@/mixins/modal';
import eventActionsMixin from '@/mixins/event-actions';

import convertObjectFieldToTreeBranch from '@/helpers/treeview';

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
      return () => this.showModal({
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
          title: this.$t('modals.createAckRemove.title'),
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },

    showVariablesHelperModal() {
      const variables = [];

      const alarmFields = convertObjectFieldToTreeBranch(omit(this.item, ['entity']), 'alarm');
      variables.push(alarmFields);

      if (this.item.entity) {
        const entityFields = convertObjectFieldToTreeBranch(this.item.entity, 'entity');
        variables.push(entityFields);
      }

      return () => this.showModal({
        name: MODALS.variablesHelp,
        config: {
          ...this.modalConfig,
          variables,
        },
      });
    },

    actionsAccessFilterHandler({ type }) {
      const right = BUSINESS_USER_RIGHTS_ACTIONS_MAP.alarmsList[type];

      if (!right) {
        return true;
      }

      return this.checkAccess(right);
    },
  },
};
