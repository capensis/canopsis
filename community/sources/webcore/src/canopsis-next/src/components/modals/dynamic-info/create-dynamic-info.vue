<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <template-testing-test-variables-wrapper
          v-model="form"
          :rule-id="dynamicInfoId"
          :type="type"
        >
          <template #default="{ templateVars, copyVars }">
            <div class="position-relative">
              <pattern-progress
                v-if="chatPending"
                :in-progress-text="chatPendingTexts.inProgress"
                :cancel-button-label="chatPendingTexts.cancel"
                @cancel="chatCancelPending"
              />
              <dynamic-info-form
                v-model="form"
                :is-disabled-id-field="isDisabledIdField"
                :template-vars="templateVars"
                :copy-vars="copyVars"
              />
            </div>
          </template>
        </template-testing-test-variables-wrapper>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled || chatPending"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DynamicInfoForm,
    PatternProgress,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.dynamicInfo;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(dynamicInfoToForm(config.value.dynamicInfo));

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancelPending: chatCancelPending,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: props.modal.id,
      ruleId: props.modal.config?.dynamicInfo?._id,
      context: LLM_SOCKET_CONTEXTS.dynamicInfos,
    });

    const dynamicInfoId = computed(() => config.value.dynamicInfo?._id);
    const title = computed(() => config.value.title || t('modals.createDynamicInfo.create.title'));
    const isDisabledIdField = computed(() => config.value.isDisabledIdField);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const data = await config.value.action?.(formToDynamicInfo(form.value));

        close();

        return data;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      dynamicInfoId,
      type,
      title,
      isDisabledIdField,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancelPending,
      submit,
      close,
    };
  },
};
</script>
