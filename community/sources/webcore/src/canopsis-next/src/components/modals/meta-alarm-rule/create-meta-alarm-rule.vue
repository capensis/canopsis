<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <v-layout class="gap-2" column>
          <c-enabled-field v-model="form.enabled" hide-details with-background />
          <template-testing-test-variables-wrapper
            v-model="form"
            :is-new="isNew"
            :type="type"
          >
            <template #default="{ templateVars }">
              <meta-alarm-rule-form
                v-model="form"
                ref="formElement"
                :active-step.sync="activeStep"
                :disabled-id-field="config.isDisabledIdField"
                :alarm-infos="alarmInfos"
                :entity-infos="entityInfos"
                :template-vars="templateVars"
              />
            </template>
          </template-testing-test-variables-wrapper>
        </v-layout>
        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
          v-on="chatOptions.on"
        />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          v-if="isLastStep"
          key="submit"
          :disabled="isDisabled || chatOptions.bind.pending"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
        <v-btn
          v-else
          key="next"
          :disabled="!isStepValid || chatOptions.bind.pending"
          type="button"
          class="primary"
          @click="next"
        >
          {{ $t('common.next') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import {
  computed,
  ref,
  onMounted,
  watch,
  nextTick,
  toRef,
} from 'vue';

import {
  LLM_SOCKET_CONTEXTS,
  META_ALARMS_FORM_STEPS,
  MODALS,
  TEMPLATE_TESTING_TEST_TYPES,
  VALIDATION_DELAY,
} from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfos } from '@/hooks/store/modules/entity-infos';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    MetaAlarmRuleForm,
    AiChatSidebar,
    TemplateTestingTestVariablesWrapper,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const { alarmInfos, entityInfos, fetchInfos } = useEntityInfos();

    const activeStep = ref(META_ALARMS_FORM_STEPS.general);
    const isStepValid = ref(false);
    const formElement = ref(null);
    const form = ref(metaAlarmRuleToForm(config.value.rule));

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.rule?._id,
      context: LLM_SOCKET_CONTEXTS.metaAlarmRule,
    });

    const isNew = computed(() => !config.value.rule?._id);
    const title = computed(() => config.value.title ?? t('modals.metaAlarmRule.create.title'));
    const isLastStep = computed(() => activeStep.value === META_ALARMS_FORM_STEPS.parameters);

    /**
     * Checks if the current step is valid based on error states
     * @returns {boolean} True if the current step has no errors, false otherwise
     */
    const isCurrentStepValid = () => {
      const { hasGeneralError, hasTypeError, hasParametersError } = formElement.value ?? {};

      return {
        [META_ALARMS_FORM_STEPS.general]: !hasGeneralError,
        [META_ALARMS_FORM_STEPS.type]: !hasTypeError,
        [META_ALARMS_FORM_STEPS.parameters]: !hasParametersError,
      }[activeStep.value];
    };

    /**
     * Validates the current step by calling the appropriate validation function
     *
     * @returns {Promise<boolean>|boolean|undefined} The validation result for the current step
     */
    const validateCurrentStepValid = () => {
      const {
        validateGeneralChildren,
        validateTypeChildren,
        validateParametersChildren,
      } = formElement.value ?? {};

      const func = {
        [META_ALARMS_FORM_STEPS.general]: validateGeneralChildren,
        [META_ALARMS_FORM_STEPS.type]: validateTypeChildren,
        [META_ALARMS_FORM_STEPS.parameters]: validateParametersChildren,
      }[activeStep.value];

      return func?.();
    };

    /**
     * Moves to the next step if the current step validation passes
     */
    const next = async () => {
      const isValid = await validateCurrentStepValid();

      if (isValid) {
        activeStep.value += 1;
      }
    };

    /**
     * Updates the step validation state by checking if the current step is valid
     */
    const handleValidationChanged = () => {
      isStepValid.value = isCurrentStepValid();
    };

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action(formToMetaAlarmRule(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useEntityInfoPropertyFetching();
    useFormConfirmableCloseModal({ form, submit, close });

    onMounted(() => {
      fetchInfos();

      handleValidationChanged();

      nextTick(() => {
        watch(() => formElement.value?.hasGeneralError, handleValidationChanged);
        watch(() => formElement.value?.hasTypeError, handleValidationChanged);
        watch(() => formElement.value?.hasParametersError, handleValidationChanged);
      });
    });

    return {
      form,
      config,
      isNew,
      type,
      title,
      activeStep,
      isStepValid,
      isLastStep,
      formElement,
      alarmInfos,
      entityInfos,
      isDisabled,
      submitting,
      chatShown,
      chatOptions,
      next,
      submit,
      close,
    };
  },
};
</script>
