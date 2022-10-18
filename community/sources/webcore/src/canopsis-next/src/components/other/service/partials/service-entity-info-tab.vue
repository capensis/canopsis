<template lang="pug">
  div
    v-layout.pa-2(v-if="availableActions.length", align-center, wrap)
      div {{ $t('common.actionsLabel') }}:
      service-entity-actions(
        :actions="availableActions",
        :assigned-instructions="entity.assigned_instructions",
        @apply="applyActionForEntity",
        @execute="executeAlarmInstruction"
      )
    service-entity-template(:entity="entity", :template="template")
</template>

<script>
import { MODALS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { getAvailableActionsByEntity } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityTemplate from './service-entity-template.vue';
import ServiceEntityActions from './service-entity-actions.vue';

export default {
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

    availableActions() {
      return getAvailableActionsByEntity(this.entity)
        .filter(this.actionsAccessFilterHandler);
    },
  },
  methods: {
    applyActionForEntity({ type }) {
      this.applyEntityAction(type, [this.entity]);
    },

    executeAlarmInstruction(assignedInstruction) {
      const refreshEntities = () => this.$emit('refresh');

      this.$modals.show({
        name: MODALS.executeRemediationInstruction,
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
