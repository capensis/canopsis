<template lang="pug">
  div
    v-layout.pa-2(justify-space-between)
      v-flex(xs11)
        v-layout(justify-space-between)
          v-layout(v-if="availableActions.length", row, align-center, wrap)
            div {{ $t('common.actionsLabel') }}:
            service-entity-actions(
              :entity="entity",
              :actions="availableActions",
              :assigned-instructions="entity.assigned_instructions",
              @apply="applyActionForEntity",
              @execute="executeAlarmInstruction"
            )
      v-tooltip(v-if="hasPbehaviors && hasAccessToManagePbehaviors", top)
        v-btn(slot="activator", small, @click="showPbehaviorsListModal")
          v-icon(small) list
        span {{ $t('modals.service.editPbehaviors') }}
    service-entity-template(:entity="entity", :template="template")
</template>

<script>
import {
  CRUD_ACTIONS,
  MODALS,
  PBEHAVIOR_TYPE_TYPES,
  USERS_PERMISSIONS,
} from '@/constants';

import { getAvailableActionsByEntity } from '@/helpers/entities/entity';
import { mapIds } from '@/helpers/entities';
import { isInstructionManual } from '@/helpers/forms/remediation-instruction';

import { authMixin } from '@/mixins/auth';
import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityTemplate from '@/components/modals/service/partial/service-entity-template.vue';
import ServiceEntityActions from '@/components/modals/service/partial/service-entity-actions.vue';

export default {
  inject: ['$actionsQueue'],
  components: {
    ServiceEntityActions,
    ServiceEntityTemplate,
  },
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
      return this.$actionsQueue.queue
        .filter(({ entities }) => mapIds(entities).includes(this.entity._id))
        .map(({ actionType }) => actionType);
    },

    availableActions() {
      return getAvailableActionsByEntity(this.entity)
        .filter(this.actionsAccessFilterHandler)
        .map(action => ({
          ...action,
          disabled: this.isActionDisabled(action.type),
        }));
    },
  },
  methods: {
    isActionDisabled(action) {
      return this.clickedActions.includes(action);
    },

    applyActionForEntity({ type }) {
      this.addEntityAction(type, [this.entity]);
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
        name: isInstructionManual(assignedInstruction)
          ? MODALS.executeRemediationInstruction
          : MODALS.executeRemediationSimpleInstruction,
        config: {
          assignedInstruction,
          alarmId: this.entity.alarm_id,
          onClose: refreshEntities,
          onComplete: refreshEntities,
        },
      });
    },
  },
};
</script>
