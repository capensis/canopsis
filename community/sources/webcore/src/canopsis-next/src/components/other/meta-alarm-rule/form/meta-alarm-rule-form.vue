<template>
  <v-stepper :value="activeStep" @change="$emit('update:active-step', $event)">
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
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        ref="typeStepElement"
        :step="META_ALARMS_FORM_STEPS.type"
        class="pa-0"
      >
        <div class="pa-4">
          <meta-alarm-rule-type-field v-field="form.type" />
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
          <meta-alarm-rule-parameters-form v-field="form" />
        </c-information-block>
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import { ref } from 'vue';

import { META_ALARMS_FORM_STEPS } from '@/constants';

import MetaAlarmRuleGeneralForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue';
import MetaAlarmRuleTypeField from '@/components/other/meta-alarm-rule/form/fields/meta-alarm-rule-type-field.vue';
import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

import { useElementChildrenValidation } from '@/hooks/form/useValidationChildren';

export default {
  components: {
    CInformationBlock,
    MetaAlarmRuleParametersForm,
    MetaAlarmRuleTypeField,
    MetaAlarmRuleGeneralForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
    activeStep: {
      type: Number,
      default: 0,
    },
  },
  setup(props, { expose }) {
    const generalStepElement = ref(null);
    const {
      hasChildrenError: hasGeneralError,
      validateChildren: validateGeneralChildren,
    } = useElementChildrenValidation(generalStepElement);

    const typeStepElement = ref(null);
    const {
      hasChildrenError: hasTypeError,
      validateChildren: validateTypeChildren,
    } = useElementChildrenValidation(typeStepElement);

    const parametersStepElement = ref(null);
    const {
      hasChildrenError: hasParametersError,
      validateChildren: validateParametersChildren,
    } = useElementChildrenValidation(parametersStepElement);

    expose({
      hasGeneralError,
      hasTypeError,
      hasParametersError,
      validateGeneralChildren,
      validateTypeChildren,
      validateParametersChildren,
    });

    return {
      generalStepElement,
      hasGeneralError,
      parametersStepElement,
      hasParametersError,
      typeStepElement,
      hasTypeError,
      META_ALARMS_FORM_STEPS,
    };
  },
};
</script>
