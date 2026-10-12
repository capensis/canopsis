<template>
  <v-layout class="gap-3" column>
    <v-layout wrap>
      <v-flex xs12>
        <v-switch
          :input-value="form.type"
          :false-value="$constants.SNMP_STATE_TYPES.simple"
          :true-value="$constants.SNMP_STATE_TYPES.template"
          :label="$t('snmpRule.toCustom')"
          color="primary"
          class="mt-0 pt-0"
          hide-details
          @change="updateTypeField"
        />
      </v-flex>
    </v-layout>
    <v-divider light />
    <template v-if="isTemplate">
      <snmp-rule-form-module-mib-objects-form
        v-field="form.stateoid"
        :items="items"
        :label="$t('snmpRule.defineVar')"
      />
      <v-layout wrap>
        <v-flex xs12>
          <v-layout
            v-for="{ value, color, whiteText, key, text } in availableStates"
            :key="value"
            class="gap-4"
            align-center
            justify-center
          >
            <v-flex xs2>
              <v-chip
                :style="{ backgroundColor: color }"
                :class="{ 'white--text': whiteText }"
                class="snmp-rule-state-chip rounded-lg"
                label
              >
                <strong class="state-title">
                  {{ text }}
                </strong>
              </v-chip>
            </v-flex>
            <v-flex xs10>
              <v-text-field
                v-field="form[key]"
                :placeholder="$t('snmpRule.writeTemplate')"
              />
            </v-flex>
          </v-layout>
        </v-flex>
      </v-layout>
    </template>
    <state-criticity-field v-else v-field="form.state" class="mt-1" />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { SNMP_STATE_TYPES, SNMP_TEMPLATE_STATE_STATES } from '@/constants';

import { getSnmpRuleStateColor } from '@/helpers/entities/snmp-rule/color';

import { useI18n } from '@/hooks/i18n';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';

export default {
  components: {
    StateCriticityField,
    SnmpRuleFormModuleMibObjectsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const isTemplate = computed(() => props.form.type === SNMP_STATE_TYPES.template);

    const availableStates = computed(() => Object.entries(props.stateValues).map(([key, state]) => ({
      key,
      text: t(`snmpRule.states.${key}`),
      value: state,
      color: getSnmpRuleStateColor(state),
      whiteText: [SNMP_TEMPLATE_STATE_STATES.info, SNMP_TEMPLATE_STATE_STATES.critical].includes(state),
    })));

    const updateTypeField = (type) => {
      const state = {
        type,
      };

      if (type === SNMP_STATE_TYPES.template) {
        state.stateoid = {};
      }

      emit('input', state);
    };

    return {
      isTemplate,
      availableStates,
      updateTypeField,
    };
  },
};
</script>

<style lang="scss" scoped>
  .snmp-rule-state-chip {
    width: 100%;

    ::v-deep .v-chip__content {
      width: 100%;
    }

    .state-title {
      text-transform: uppercase;
      margin: auto;
    }
  }
</style>
