<template>
  <v-stepper v-model="activeStepComputed">
    <v-stepper-header>
      <v-stepper-step
        :complete="activeStep > META_ALARMS_FORM_STEPS.general"
        :step="META_ALARMS_FORM_STEPS.general"
        :rules="[() => !hasGeneralError]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.basics') }}
        <small v-if="hasGeneralError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="activeStep > META_ALARMS_FORM_STEPS.type"
        :step="META_ALARMS_FORM_STEPS.type"
        :rules="[() => !hasTypeError]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.defineType') }}
        <small v-if="hasTypeError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="activeStep > META_ALARMS_FORM_STEPS.parameters"
        :step="META_ALARMS_FORM_STEPS.parameters"
        :rules="[() => !hasParametersError]"
        class="py-0"
        editable
      >
        {{ $t('metaAlarmRule.steps.addParameters') }}
        <small v-if="hasParametersError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
    </v-stepper-header>

    <v-stepper-items>
      <v-stepper-content
        ref="generalStepElement"
        :step="META_ALARMS_FORM_STEPS.general"
        class="pa-0"
      >
        <meta-alarm-rule-general-form
          v-field="form"
          :disabled-id-field="disabledIdField"
          :variables="templateVars.output"
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        ref="typeStepElement"
        :step="META_ALARMS_FORM_STEPS.type"
        class="pa-0"
      >
        <div class="pa-4">
          <meta-alarm-rule-type-form
            v-field="form"
            :variables="templateVars.entity"
          />
        </div>
      </v-stepper-content>
      <v-stepper-content
        ref="parametersStepElement"
        :step="META_ALARMS_FORM_STEPS.parameters"
        class="pa-0"
      >
        <c-information-block
          :title="$t(`metaAlarmRule.parametersTitle.${form.type}`)"
          class="pa-4"
        >
          <span class="text--secondary mb-2">{{ $t(`metaAlarmRule.parametersDescription.${form.type}`) }}</span>
          <meta-alarm-rule-parameters-form v-field="form" :template-vars="templateVars" />
        </c-information-block>
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import { computed, ref } from 'vue';

import { META_ALARMS_FORM_STEPS, META_ALARMS_RULE_TYPES } from '@/constants';

import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';
import MetaAlarmRuleTypeForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-type-form.vue';
import MetaAlarmRuleGeneralForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue';

export default {
  components: {
    MetaAlarmRuleParametersForm,
    MetaAlarmRuleTypeForm,
    MetaAlarmRuleGeneralForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({
        type: META_ALARMS_RULE_TYPES.attribute,
      }),
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
    activeStep: {
      type: Number,
      default: META_ALARMS_FORM_STEPS.general,
    },
    alarmInfos: {
      type: Array,
      default: () => [],
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { expose, emit }) {
    const generalStepElement = ref(null);
    const typeStepElement = ref(null);
    const parametersStepElement = ref(null);

    const {
      hasChildrenError: hasGeneralError,
      validateChildren: validateGeneralChildren,
    } = useValidationElementChildren(generalStepElement);

    const {
      hasChildrenError: hasTypeError,
      validateChildren: validateTypeChildren,
    } = useValidationElementChildren(typeStepElement);

    const {
      hasChildrenError: hasParametersError,
      validateChildren: validateParametersChildren,
    } = useValidationElementChildren(parametersStepElement);

    const activeStepComputed = computed({
      get: () => props.activeStep,
      set: value => emit('update:active-step', value),
    });

    useAiChatExpand({ activeTab: activeStepComputed, neededTab: META_ALARMS_FORM_STEPS.parameters });

    expose({
      hasGeneralError,
      hasTypeError,
      hasParametersError,
      validateGeneralChildren,
      validateTypeChildren,
      validateParametersChildren,
    });

    return {
      META_ALARMS_FORM_STEPS,

      activeStepComputed,

      generalStepElement,
      typeStepElement,
      parametersStepElement,

      hasGeneralError,
      hasParametersError,
      hasTypeError,
    };
  },
};
</script>
