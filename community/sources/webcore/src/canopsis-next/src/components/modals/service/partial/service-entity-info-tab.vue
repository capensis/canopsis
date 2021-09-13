<template lang="pug">
  div
    v-layout.pa-2(justify-space-between)
      v-flex(xs11)
        v-layout(justify-space-between)
          v-layout(v-if="availableActions.length", row, align-center)
            div {{ $t('common.actionsLabel') }}:
            div(v-for="action in availableActions", :key="action.eventType")
              v-tooltip(top)
                v-btn(
                  slot="activator",
                  :disabled="isActionDisabled(action.eventType)",
                  depressed,
                  small,
                  light,
                  @click.stop="action.action"
                )
                  v-icon {{ action.icon }}
                span {{ $t(`common.actions.${action.eventType}`) }}
      v-tooltip(v-if="active && hasAccessToManagePbehaviors", top)
        v-btn(slot="activator", small, @click="showPbehaviorsListModal")
          v-icon(small) list
        span {{ $t('modals.service.editPbehaviors') }}
    service-entity-template(:entity="entity", :template="template")
</template>

<script>
import { isNull, pickBy } from 'lodash';

import {
  CRUD_ACTIONS,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  MODALS,
  PBEHAVIOR_TYPE_TYPES,
  USERS_PERMISSIONS,
  WIDGETS_ACTIONS_TYPES,
} from '@/constants';

import { authMixin } from '@/mixins/auth';
import widgetActionPanelServiceEntityMixin from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityTemplate from '@/components/modals/service/partial/service-entity-template.vue';

export default {
  inject: ['$eventsQueue'],
  components: { ServiceEntityTemplate },
  mixins: [
    authMixin,
    widgetActionPanelServiceEntityMixin,
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      default: '',
    },
    active: {
      type: Boolean,
      default: false,
    },
    paused: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    pausedPbehaviors() {
      return this.entity.pbehaviors.filter(pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause);
    },

    hasAccessToManagePbehaviors() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors);
    },

    clickedActions() {
      return this.$eventsQueue.queue
        .reduce((acc, { entityId, type }) => {
          if (entityId === this.entity._id) {
            acc.push(type);
          }

          return acc;
        }, []);
    },

    actionsMap() {
      const { weather: weatherActionsTypes } = WIDGETS_ACTIONS_TYPES;

      return {
        ack: {
          type: weatherActionsTypes.entityAck,
          eventType: EVENT_ENTITY_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          action: this.prepareAckAction,
        },
        assocTicket: {
          type: weatherActionsTypes.entityAssocTicket,
          eventType: EVENT_ENTITY_TYPES.assocTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          action: this.prepareAssocTicketAction,
        },
        validate: {
          type: weatherActionsTypes.entityValidate,
          eventType: EVENT_ENTITY_TYPES.validate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.validate].icon,
          action: this.prepareValidateAction,
        },
        invalidate: {
          type: weatherActionsTypes.entityInvalidate,
          eventType: EVENT_ENTITY_TYPES.invalidate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.invalidate].icon,
          action: this.prepareInvalidateAction,
        },
        pause: {
          type: weatherActionsTypes.entityPause,
          eventType: EVENT_ENTITY_TYPES.pause,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pause].icon,
          action: this.preparePauseAction,
        },
        play: {
          type: weatherActionsTypes.entityPlay,
          eventType: EVENT_ENTITY_TYPES.play,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.play].icon,
          action: this.preparePlayAction,
        },
        cancel: {
          type: weatherActionsTypes.entityCancel,
          eventType: EVENT_ENTITY_TYPES.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          action: this.prepareCancelAction,
        },
        comment: {
          type: weatherActionsTypes.entityComment,
          eventType: EVENT_ENTITY_TYPES.comment,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon,
          action: this.prepareCommentAction,
        },
      };
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    availableActions() {
      const { filteredActionsMap } = this;
      const actions = [filteredActionsMap.comment];

      if (this.entity.state.val !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
        actions.push(filteredActionsMap.ack);
      }

      actions.push(filteredActionsMap.assocTicket);

      if (this.entity.state.val === ENTITIES_STATES.major) {
        actions.push(filteredActionsMap.validate, filteredActionsMap.invalidate);
      }

      if (this.paused) {
        actions.push(filteredActionsMap.play);
      } else {
        actions.push(filteredActionsMap.pause);
      }

      if (
        this.entity.alarm_display_name
        && (!this.entity.status || this.entity.status.val !== ENTITIES_STATUSES.cancelled)
      ) {
        actions.push(filteredActionsMap.cancel);
      }

      return actions.filter(action => !!action);
    },
  },
  methods: {
    isActionDisabled(action) {
      const SERVICE_ENTITY_EVENT_TYPES_TO_EVENT_ENTITY_TYPES = {
        [EVENT_ENTITY_TYPES.validate]: EVENT_ENTITY_TYPES.changeState,
        [EVENT_ENTITY_TYPES.invalidate]: EVENT_ENTITY_TYPES.cancel,
      };

      return this.clickedActions.includes(SERVICE_ENTITY_EVENT_TYPES_TO_EVENT_ENTITY_TYPES[action] || action);
    },

    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.pausedPbehaviors,
          entityId: this.entity._id,
          onlyActive: true,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
    },
  },
};
</script>
