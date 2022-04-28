<template lang="pug">
  div
    v-layout.pa-2(justify-space-between)
      v-flex(xs11)
        v-layout(justify-space-between)
          v-layout(v-if="availableActions.length", row, align-center)
            div {{ $t('common.actionsLabel') }}:
            div(v-for="action in availableActions", :key="action.eventType")
              v-tooltip(top)
                template(slot="activator")
                  service-entity-alarm-instruction-menu(
                    v-if="action.eventType === $constants.EVENT_ENTITY_TYPES.executeInstruction",
                    :assigned-instructions="entity.assigned_instructions",
                    :icon="action.icon",
                    @execute="executeAlarmInstruction"
                  )
                  v-btn(
                    v-else,
                    :disabled="isActionDisabled(action.eventType)",
                    depressed,
                    small,
                    light,
                    @click.stop="action.action"
                  )
                    v-icon {{ action.icon }}
                span {{ $t(`common.actions.${action.eventType}`) }}
      v-tooltip(v-if="hasPbehaviors && hasAccessToManagePbehaviors", top)
        v-btn(slot="activator", small, @click="showPbehaviorsListModal")
          v-icon(small) list
        span {{ $t('modals.service.editPbehaviors') }}
    service-entity-template(:entity="entity", :template="template")
</template>

<script>
import { isNull } from 'lodash';

import {
  CRUD_ACTIONS,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  MODALS, PBEHAVIOR_TYPE_TYPES,
  USERS_PERMISSIONS,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import { authMixin } from '@/mixins/auth';
import widgetActionPanelServiceEntityMixin from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityTemplate from '@/components/modals/service/partial/service-entity-template.vue';
import ServiceEntityAlarmInstructionMenu from './service-entity-alarm-instruction-menu.vue';

export default {
  inject: ['$eventsQueue'],
  components: { ServiceEntityAlarmInstructionMenu, ServiceEntityTemplate },
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
  },
  computed: {
    paused() {
      return this.entity.pbehavior_info?.canonical_type === PBEHAVIOR_TYPE_TYPES.pause;
    },

    hasPbehaviors() {
      return !!this.entity.pbehaviors.length;
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
      return {
        ack: {
          type: WEATHER_ACTIONS_TYPES.entityAck,
          eventType: EVENT_ENTITY_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          action: this.prepareAckAction,
        },
        assocTicket: {
          type: WEATHER_ACTIONS_TYPES.entityAssocTicket,
          eventType: EVENT_ENTITY_TYPES.assocTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          action: this.prepareAssocTicketAction,
        },
        executeInstruction: {
          type: WEATHER_ACTIONS_TYPES.executeInstruction,
          eventType: EVENT_ENTITY_TYPES.executeInstruction,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.executeInstruction].icon,
        },
        validate: {
          type: WEATHER_ACTIONS_TYPES.entityValidate,
          eventType: EVENT_ENTITY_TYPES.validate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.validate].icon,
          action: this.prepareValidateAction,
        },
        invalidate: {
          type: WEATHER_ACTIONS_TYPES.entityInvalidate,
          eventType: EVENT_ENTITY_TYPES.invalidate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.invalidate].icon,
          action: this.prepareInvalidateAction,
        },
        pause: {
          type: WEATHER_ACTIONS_TYPES.entityPause,
          eventType: EVENT_ENTITY_TYPES.pause,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pause].icon,
          action: this.preparePauseAction,
        },
        play: {
          type: WEATHER_ACTIONS_TYPES.entityPlay,
          eventType: EVENT_ENTITY_TYPES.play,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.play].icon,
          action: this.preparePlayAction,
        },
        cancel: {
          type: WEATHER_ACTIONS_TYPES.entityCancel,
          eventType: EVENT_ENTITY_TYPES.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          action: this.prepareCancelAction,
        },
        comment: {
          type: WEATHER_ACTIONS_TYPES.entityComment,
          eventType: EVENT_ENTITY_TYPES.comment,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon,
          action: this.prepareCommentAction,
        },
      };
    },

    availableActions() {
      const { actionsMap } = this;
      const actions = [actionsMap.comment];

      if (this.entity.assigned_instructions && this.entity.assigned_instructions.length) {
        actions.push(actionsMap.executeInstruction);
      }

      if (this.entity.state.val !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
        actions.push(actionsMap.ack);
      }

      actions.push(actionsMap.assocTicket);

      if (this.entity.state.val === ENTITIES_STATES.major) {
        actions.push(actionsMap.validate, actionsMap.invalidate);
      }

      if (this.paused) {
        actions.push(actionsMap.play);
      } else {
        actions.push(actionsMap.pause);
      }

      if (
        this.entity.alarm_display_name
        && (!this.entity.status || this.entity.status.val !== ENTITIES_STATUSES.cancelled)
      ) {
        actions.push(actionsMap.cancel);
      }

      return actions.filter(this.actionsAccessFilterHandler);
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
          pbehaviors: this.entity.pbehaviors,
          entityId: this.entity._id,
          onlyActive: true,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
    },

    executeAlarmInstruction(assignedInstruction) {
      const refreshEntities = () => this.$emit('refresh');

      this.$modals.show({
        name: MODALS.executeRemediationInstruction,
        config: {
          assignedInstruction,
          alarmId: this.entity.alarm_id,
          onOpen: refreshEntities,
          onClose: refreshEntities,
          onComplete: refreshEntities,
        },
      });
    },
  },
};
</script>
