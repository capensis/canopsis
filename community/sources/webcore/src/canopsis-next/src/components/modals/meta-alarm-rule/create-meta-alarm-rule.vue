<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <template-testing-test-variables-wrapper
        v-field="form"
        :is-new="isNew"
        :type="type"
      >
        <template #default="{ templateVars }">
          <v-form>
            <meta-alarm-rule-form
              v-model="form"
              ref="formElement"
              :active-step.sync="activeStep"
              :disabled-id-field="config.isDisabledIdField"
              :alarm-infos="alarmInfos"
              :entity-infos="entityInfos"
              :template-vars="templateVars"
            />
          </v-form>
        </template>
      </template-testing-test-variables-wrapper>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        v-if="isLastStep"
        key="submit"
        :disabled="isDisabled"
        :loading="submitting"
        class="primary"
        type="submit"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
      <v-btn
        v-else
        key="next"
        :disabled="!isStepValid"
        type="button"
        class="primary"
        @click="next"
      >
        {{ $t('common.next') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import {
  computed,
  ref,
  onMounted,
  watch,
  nextTick,
} from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY, META_ALARMS_FORM_STEPS } from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useEntityInfos } from '@/hooks/store/modules/entity-infos';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    MetaAlarmRuleForm,
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

    const isNew = computed(() => !config.value.rule?._id);
    const title = computed(() => config.value.title ?? t('modals.metaAlarmRule.create.title'));
    const isLastStep = computed(() => activeStep.value === META_ALARMS_FORM_STEPS.parameters);

    const isCurrentStepValid = () => {
      const { hasGeneralError, hasTypeError, hasParametersError } = formElement.value ?? {};

      return {
        [META_ALARMS_FORM_STEPS.general]: !hasGeneralError,
        [META_ALARMS_FORM_STEPS.type]: !hasTypeError,
        [META_ALARMS_FORM_STEPS.parameters]: !hasParametersError,
      }[activeStep.value];
    };

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

    const next = async () => {
      const isValid = await validateCurrentStepValid();

      if (isValid) {
        activeStep.value += 1;
      }
    };

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action(formToMetaAlarmRule(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    onMounted(() => {
      fetchInfos();

      const handleValidationChanged = () => {
        isStepValid.value = isCurrentStepValid();
      };

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
      next,
      submit,
    };
  },
};
</script>
