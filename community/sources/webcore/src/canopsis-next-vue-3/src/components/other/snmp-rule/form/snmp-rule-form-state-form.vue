<template>
  <c-layout gap="3" column>
    <v-row density="comfortable">
      <v-col cols="12">
        <v-switch
          :model-value="form.type"
          :false-value="SNMP_STATE_TYPES.simple"
          :true-value="SNMP_STATE_TYPES.template"
          :label="t('snmpRule.toCustom')"
          color="primary"
          class="mt-0 pt-0"
          hide-details
          @update:model-value="updateTypeField"
        />
      </v-col>
    </v-row>

    <v-divider />

    <template v-if="isTemplate">
      <snmp-rule-form-module-mib-objects-form
        :form="form.stateoid"
        :items="items"
        :label="t('snmpRule.defineVar')"
        name="state.stateoid"
        @update:form="updateField('stateoid', $event)"
      />

      <v-row density="comfortable">
        <v-col cols="12">
          <v-row
            v-for="{ value, color, whiteText, key, text } in availableStates"
            :key="value"
            class="align-center justify-center"
            density="comfortable"
          >
            <v-col cols="2">
              <v-chip
                :style="{ backgroundColor: color }"
                :class="{ 'text-white': whiteText }"
                class="snmp-rule-state-chip rounded-lg"
                label
              >
                <strong class="state-title">
                  {{ text }}
                </strong>
              </v-chip>
            </v-col>
            <v-col cols="10">
              <c-text-field
                :model-value="form[key]"
                :placeholder="t('snmpRule.writeTemplate')"
                :name="`state.${key}`"
                @update:model-value="updateField(key, $event)"
              />
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </template>

    <state-criticity-field
      v-else
      :value="form.state"
      class="mt-1"
      @update:value="updateField('state', $event)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { SNMP_STATE_TYPES, SNMP_TEMPLATE_STATE_STATES } from '@/constants';

import { getSnmpRuleStateColor } from '@/helpers/entities/snmp-rule/color';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';

export default {
  components: {
    StateCriticityField,
    SnmpRuleFormModuleMibObjectsForm,
  },

  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    stateValues: {
      type: Object,
      default: () => SNMP_TEMPLATE_STATE_STATES,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel, updateField } = useModelField(props, emit);

    const isTemplate = computed(() => props.form.type === SNMP_STATE_TYPES.template);

    const availableStates = computed(() =>
      Object.entries(props.stateValues).map(([key, state]) => ({
        key,
        text: t(`snmpRule.states.${key}`),
        value: state,
        color: getSnmpRuleStateColor(state),
        whiteText: [SNMP_TEMPLATE_STATE_STATES.info, SNMP_TEMPLATE_STATE_STATES.critical].includes(
          state,
        ),
      })),
    );

    const updateTypeField = (type) => {
      const state = {
        type,
      };

      if (type === SNMP_STATE_TYPES.template) {
        state.stateoid = {};
      }

      updateModel(state);
    };

    return {
      t,
      SNMP_STATE_TYPES,
      isTemplate,
      availableStates,
      updateField,
      updateTypeField,
    };
  },
};
</script>

<style lang="scss" scoped>
.snmp-rule-state-chip {
  width: 100%;

  :deep(.v-chip__content) {
    width: 100%;
  }

  .state-title {
    text-transform: uppercase;
    margin: auto;
  }
}
</style>
