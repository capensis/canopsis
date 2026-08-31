<template>
  <c-layout column>
    <c-name-field
      v-model="nameField"
      :max-length="255"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('common.priority')">
        <c-priority-field v-model="priorityField" />
      </c-form-block-row>

      <c-form-block-row v-if="!isEntityType" :label="t('common.type')" indented>
        <idle-rule-alarm-type-field v-model="alarmConditionField" :label="t('common.type')" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.description')">
        <c-description-field v-model="descriptionField" required />
      </c-form-block-row>

      <c-form-block-row :label="t('idleRules.timeRangeAwaiting')">
        <c-duration-field
          v-model:duration="durationField"
          :label="t('idleRules.timeRangeAwaiting')"
          name="duration"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.disableDuringPeriods')">
        <c-disable-during-periods-field v-model="disableDuringPeriodsField" />
      </c-form-block-row>

      <template v-if="!isEntityType">
        <c-form-block-row :label="t('common.type')">
          <c-action-type-field
            v-model="operationTypeField"
            :types="actionTypes"
            name="operation.type"
          />
        </c-form-block-row>

        <action-parameters-form
          v-model="parameters"
          :type="operationTypeField"
          :depth="1"
          name="operation.parameters"
        />

        <c-form-block-row
          v-if="isAssociateTicketAction"
          :label="t('common.comment')"
          :depth="1"
          top-border
        >
          <c-description-field v-model="commentField" :label="t('common.comment')" name="comment" />
        </c-form-block-row>
      </template>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { ACTION_TYPES } from '@/constants';

import { isAssociateTicketActionType } from '@/helpers/entities/action';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';

import IdleRuleAlarmTypeField from './fields/idle-rule-alarm-type-field.vue';

const IDLE_RULE_ACTION_TYPES = [
  ACTION_TYPES.snooze,
  ACTION_TYPES.ack,
  ACTION_TYPES.ackremove,
  ACTION_TYPES.cancel,
  ACTION_TYPES.assocticket,
  ACTION_TYPES.changeState,
  ACTION_TYPES.pbehavior,
];

export default {
  components: {
    ActionParametersForm,
    IdleRuleAlarmTypeField,
  },

  props: {
    isEntityType: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();

    const {
      nameField,
      priorityField,
      alarmConditionField,
      descriptionField,
      durationField,
      disableDuringPeriodsField,
      operationTypeField,
      commentField,
    } = form.models({
      nameField: 'name',
      priorityField: 'priority',
      alarmConditionField: 'alarm_condition',
      descriptionField: 'description',
      durationField: 'duration',
      disableDuringPeriodsField: 'disable_during_periods',
      operationTypeField: 'operation.type',
      commentField: 'comment',
    });

    const parameters = computed({
      get: () => form.getValue(`operation.parameters.${operationTypeField.value}`),
      set: value => form.setValue(`operation.parameters.${operationTypeField.value}`, value),
    });

    const isAssociateTicketAction = computed(() => isAssociateTicketActionType(operationTypeField.value),
    );

    return {
      t,
      nameField,
      priorityField,
      alarmConditionField,
      descriptionField,
      durationField,
      disableDuringPeriodsField,
      operationTypeField,
      commentField,
      parameters,
      actionTypes: IDLE_RULE_ACTION_TYPES,
      isAssociateTicketAction,
    };
  },
};
</script>
