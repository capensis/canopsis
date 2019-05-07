import { omit } from 'lodash';

import { MODALS, EVENT_ENTITY_TYPES, BUSINESS_USER_RIGHTS_ACTIONS_MAP, CRUD_ACTIONS } from '@/constants';

import modalMixin from '@/mixins/modal';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import { convertObjectToTreeview } from '@/helpers/treeview';

/**
 * @mixin Mixin for the alarms list actions panel, show modal of the action
 */
export default {
  mixins: [modalMixin, eventActionsAlarmMixin, entitiesPbehaviorMixin],
  methods: {
    showActionModal(name) {
      return () => this.showModal({
        name,
        config: this.modalConfig,
      });
    },

    showAckModal() {
      this.showModal({
        name: MODALS.createAckEvent,
        config: {
          ...this.modalConfig,
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
        },
      });
    },

    showPbehaviorsListModal() {
      this.showModal({
        name: MODALS.pbehaviorList,
        config: {
          ...this.modalConfig,
          pbehaviors: this.item.pbehaviors,
          entityId: this.item.entity._id,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
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
          title: this.$t('modals.createAckRemove.title'),
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },

    showVariablesHelperModal() {
      const variables = [];

      const alarmFields = convertObjectToTreeview(omit(this.item, ['entity']), 'alarm');

      variables.push(alarmFields);

      if (this.item.entity) {
        const entityFields = convertObjectToTreeview(this.item.entity, 'entity');
        variables.push(entityFields);
      }

      this.showModal({
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
