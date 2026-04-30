<template>
  <v-stepper
    v-model="stepper"
    non-linear
  >
    <v-stepper-header>
      <v-stepper-step
        :complete="stepper > DYNAMIC_INFO_FORM_STEPS.general"
        :step="DYNAMIC_INFO_FORM_STEPS.general"
        :rules="[() => !hasGeneralFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('common.general') }}
        <small v-if="hasGeneralFormAnyError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > DYNAMIC_INFO_FORM_STEPS.infos"
        :step="DYNAMIC_INFO_FORM_STEPS.infos"
        :rules="[() => !hasInfosFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.infos.title') }}
        <small v-if="hasInfosFormAnyError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
      <v-divider />
      <v-stepper-step
        :complete="stepper > DYNAMIC_INFO_FORM_STEPS.patterns"
        :step="DYNAMIC_INFO_FORM_STEPS.patterns"
        :rules="[() => !hasPatternsFormAnyError]"
        class="py-0"
        editable
      >
        {{ $t('modals.createDynamicInfo.steps.patterns.title') }}
        <small v-if="hasPatternsFormAnyError">{{ $t('errors.invalid') }}</small>
      </v-stepper-step>
    </v-stepper-header>
    <v-stepper-items>
      <v-stepper-content
        :step="DYNAMIC_INFO_FORM_STEPS.general"
        class="pa-0"
      >
        <dynamic-info-general-form
          v-field="form"
          ref="generalFormElement"
          :is-disabled-id-field="isDisabledIdField"
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        :step="DYNAMIC_INFO_FORM_STEPS.infos"
        class="pa-0"
      >
        <dynamic-info-infos-form
          v-field="form.infos"
          ref="infosFormElement"
          :variables="templateVars.value"
          :copy-variables="copyVars.value"
          class="pa-4"
        />
      </v-stepper-content>
      <v-stepper-content
        :step="DYNAMIC_INFO_FORM_STEPS.patterns"
        class="pa-0"
      >
        <dynamic-info-patterns-form
          v-field="form.patterns"
          ref="patternsFormElement"
          class="pa-4"
        />
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
import { inject, ref, watch } from 'vue';

import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

import DynamicInfoGeneralForm from './fields/dynamic-info-general-form.vue';
import DynamicInfoInfosForm from './fields/dynamic-info-infos-form.vue';
import DynamicInfoPatternsForm from './fields/dynamic-info-patterns-form.vue';

const DYNAMIC_INFO_FORM_STEPS = {
  general: 1,
  infos: 2,
  patterns: 3,
};

export default {
  components: {
    DynamicInfoGeneralForm,
    DynamicInfoInfosForm,
    DynamicInfoPatternsForm,
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
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    copyVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    inject('$validator');

    const stepper = ref(1);
    const hasGeneralFormAnyError = ref(false);
    const hasInfosFormAnyError = ref(false);
    const hasPatternsFormAnyError = ref(false);

    const generalFormElement = ref(null);
    const infosFormElement = ref(null);
    const patternsFormElement = ref(null);

    useAiChatExpand({ activeTab: stepper, neededTab: DYNAMIC_INFO_FORM_STEPS.patterns });

    watch(() => generalFormElement.value?.hasAnyError, (value) => {
      hasGeneralFormAnyError.value = value ?? false;
    });

    watch(() => infosFormElement.value?.hasAnyError, (value) => {
      hasInfosFormAnyError.value = value ?? false;
    });

    watch(() => patternsFormElement.value?.hasAnyError, (value) => {
      hasPatternsFormAnyError.value = value ?? false;
    });

    return {
      DYNAMIC_INFO_FORM_STEPS,
      stepper,
      hasGeneralFormAnyError,
      hasInfosFormAnyError,
      hasPatternsFormAnyError,
      generalFormElement,
      infosFormElement,
      patternsFormElement,
    };
  },
};
</script>
