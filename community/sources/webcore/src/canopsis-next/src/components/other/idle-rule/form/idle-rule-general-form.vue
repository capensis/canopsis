<template>
  <v-layout column>
    <template v-if="!isEntityType">
      <idle-rule-alarm-type-field
        class="mb-2"
        v-field="form.alarm_condition"
        :label="$t('common.type')"
      />
    </template>
    <c-name-field
      v-field="form.name"
      required
    />
    <c-description-field
      v-field="form.description"
      required
    />
    <v-layout justify-space-between>
      <v-flex xs7>
        <c-duration-field
          v-field="form.duration"
          :label="$t('idleRules.timeRangeAwaiting')"
          required
        />
      </v-flex>
      <v-flex xs3>
        <c-priority-field v-field="form.priority" />
      </v-flex>
    </v-layout>
    <c-disable-during-periods-field v-field="form.disable_during_periods" />
    <template v-if="!isEntityType">
      <c-action-type-field
        v-field="form.operation.type"
        :types="actionTypes"
        name="operation.type"
      />
      <action-parameters-form
        v-model="parameters"
        :type="form.operation.type"
        name="operation.parameters"
      />
      <c-description-field
        v-if="isAssociateTicketAction"
        v-field="form.comment"
        :label="$tc('common.comment')"
        name="comment"
      />
    </template>
  </v-layout>
</template>

<script>
import { ACTION_TYPES } from '@/constants';

import { isAssociateTicketActionType } from '@/helpers/entities/action';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import ActionParametersForm from '@/components/other/action/form/action-parameters-form.vue';

import IdleRuleAlarmTypeField from './fields/idle-rule-alarm-type-field.vue';

export default {
  inject: ['$validator'],
  components: { IdleRuleAlarmTypeField, ActionParametersForm },
  mixins: [formMixin, formValidationHeaderMixin],
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
  computed: {
    parameters: {
      get() {
        const { type, parameters } = this.form.operation;

        return parameters[type];
      },

      set(value) {
        this.updateField(`operation.parameters.${this.form.operation.type}`, value);
      },
    },

    isAssociateTicketAction() {
      return isAssociateTicketActionType(this.form.operation.type);
    },

    actionTypes() {
      return [
        ACTION_TYPES.snooze,
        ACTION_TYPES.ack,
        ACTION_TYPES.ackremove,
        ACTION_TYPES.cancel,
        ACTION_TYPES.assocticket,
        ACTION_TYPES.changeState,
        ACTION_TYPES.pbehavior,
      ];
    },
  },
};
</script>
