<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :max-length="255"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field v-field="form.priority" />
      </c-form-block-row>

      <c-form-block-row v-if="!isEntityType" :label="$t('common.type')" indented>
        <idle-rule-alarm-type-field
          v-field="form.alarm_condition"
          :label="$t('common.type')"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field
          v-field="form.description"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('idleRules.timeRangeAwaiting')">
        <c-duration-field
          v-field="form.duration"
          :label="$t('idleRules.timeRangeAwaiting')"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.disableDuringPeriods')">
        <c-disable-during-periods-field v-field="form.disable_during_periods" />
      </c-form-block-row>

      <template v-if="!isEntityType">
        <c-form-block-row :label="$t('common.type')">
          <c-action-type-field
            v-field="form.operation.type"
            :types="actionTypes"
            name="operation.type"
          />
        </c-form-block-row>

        <action-parameters-form
          v-model="parameters"
          :type="form.operation.type"
          :depth="1"
          name="operation.parameters"
        />

        <c-form-block-row
          v-if="isAssociateTicketAction"
          :label="$tc('common.comment')"
          :depth="1"
          top-border
        >
          <c-description-field
            v-field="form.comment"
            :label="$tc('common.comment')"
            name="comment"
          />
        </c-form-block-row>
      </template>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { ACTION_TYPES } from '@/constants';

import { isAssociateTicketActionType } from '@/helpers/entities/action';

import { useModelField } from '@/hooks/form/model-field';

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
  inject: ['$validator'],
  components: { IdleRuleAlarmTypeField, ActionParametersForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isEntityType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const parameters = computed({
      get() {
        const { type, parameters: operationParameters } = props.form.operation;

        return operationParameters[type];
      },

      set(value) {
        updateField(`operation.parameters.${props.form.operation.type}`, value);
      },
    });

    const isAssociateTicketAction = computed(() => (
      isAssociateTicketActionType(props.form.operation.type)
    ));

    return {
      parameters,
      actionTypes: IDLE_RULE_ACTION_TYPES,
      isAssociateTicketAction,
    };
  },
};
</script>
