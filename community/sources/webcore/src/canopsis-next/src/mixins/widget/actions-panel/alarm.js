import { get, cloneDeep } from 'lodash';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  CRUD_ACTIONS,
  QUICK_RANGES,
} from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

/**
 * @mixin Mixin for the alarms list actions panel, show modal of the action
 */
export const widgetActionsPanelAlarmMixin = {
  mixins: [
    authMixin,
    queryMixin,
    eventActionsAlarmMixin,
    entitiesPbehaviorMixin,
  ],
  methods: {
    createFastAckEvent() {
      let eventData = {};

      if (this.widget.parameters.fastAckOutput && this.widget.parameters.fastAckOutput.enabled) {
        eventData = { output: this.widget.parameters.fastAckOutput.value };
      }

      return this.createEvent(EVENT_ENTITY_TYPES.ack, this.item, eventData);
    },

    showCreateCommentModal() {
      this.$modals.show({
        name: MODALS.createCommentEvent,
        config: {
          ...this.modalConfig,
          action: data => this.createEvent(EVENT_ENTITY_TYPES.comment, this.item, data),
        },
      });
    },

    showActionModal(name) {
      return () => this.$modals.show({
        name,
        config: this.modalConfig,
      });
    },

    showSnoozeModal() {
      this.$modals.show({
        name: MODALS.createSnoozeEvent,
        config: {
          ...this.modalConfig,
          isNoteRequired: this.widget.parameters.isSnoozeNoteRequired,
        },
      });
    },

    showAckModal() {
      this.$modals.show({
        name: MODALS.createAckEvent,
        config: {
          ...this.modalConfig,
          isNoteRequired: this.widget.parameters.isAckNoteRequired,
        },
      });
    },

    showPbehaviorsListModal() {
      const availableActions = !this.isResolvedAlarm ? [CRUD_ACTIONS.delete, CRUD_ACTIONS.update] : [];

      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          ...this.modalConfig,
          pbehaviors: [this.item.pbehavior],
          entityId: this.item.entity._id,
          availableActions,
        },
      });
    },

    showCancelEventModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,
          title: this.$t('modals.createCancelEvent.title'),
          eventType: EVENT_ENTITY_TYPES.cancel,
        },
      });
    },

    showAckRemoveModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,
          title: this.$t('modals.createAckRemove.title'),
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },

    showVariablesHelperModal() {
      const {
        entity,
        pbehavior,
        infos,
        ...alarm
      } = this.item;
      const variables = [];

      variables.push(convertObjectToTreeview(alarm, 'alarm'));

      if (entity) {
        variables.push(convertObjectToTreeview(entity, 'entity'));
      }

      if (pbehavior) {
        variables.push(convertObjectToTreeview(pbehavior, 'pbehavior'));
      }

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          ...this.modalConfig,

          variables,
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: [this.item.entity._id] },
          },
        },
      });
    },

    showHistoryModal() {
      const widget = generateDefaultAlarmListWidget();

      const filter = { $and: [{ 'entity._id': get(this.item, 'entity._id') }] };
      const entityFilter = {
        title: this.item.entity.name,
        filter,
      };

      /**
       * Default value for columns
       */
      widget.parameters.widgetColumns = cloneDeep(this.widget.parameters.widgetColumns);

      /**
       * Default value for liveReporting is last 30 days
       */
      widget.parameters.liveReporting = {
        tstart: QUICK_RANGES.last30Days.start,
        tstop: QUICK_RANGES.last30Days.stop,
      };

      /**
       * Default value for opened
       */
      widget.parameters.opened = false;

      /**
       * Special entity filter for alarms list modal
       */
      widget.parameters.mainFilter = entityFilter;
      widget.parameters.viewFilters = [entityFilter];

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
        },
      });
    },

    showManualMetaAlarmUngroupModal() {
      this.$modals.show({
        name: MODALS.createEvent,
        config: {
          ...this.modalConfig,
          title: this.$t('alarmList.actions.titles.manualMetaAlarmUngroup'),
          eventType: EVENT_ENTITY_TYPES.manualMetaAlarmUngroup,
          parentsIds: [get(this.parentAlarm, 'd')],
        },
      });
    },

    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },
  },
};
