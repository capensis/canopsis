<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-2" column>
          <c-enabled-field v-model="form.enabled" with-background />
          <c-form-general-patterns-tabs
            v-model="form"
            :rule-id="dynamicInfoId"
            :type="type"
          >
            <template #general="{ setRef, templateVars, copyVars }">
              <dynamic-info-general-form
                v-model="form"
                :ref="setRef"
                :copy-vars="copyVars"
                :is-disabled-id-field="isDisabledIdField"
                :template-vars="templateVars"
              />
            </template>
            <template #patterns="{ setRef }">
              <dynamic-info-patterns-form
                v-model="form.patterns"
                :ref="setRef"
              />
            </template>
          </c-form-general-patterns-tabs>
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
          :disabled="chatOptions.bind.pending"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import DynamicInfoGeneralForm from '@/components/other/dynamic-info/form/dynamic-info-general-form.vue';
import DynamicInfoPatternsForm from '@/components/other/dynamic-info/form/dynamic-info-patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    AiChatSidebar,
    DynamicInfoGeneralForm,
    DynamicInfoPatternsForm,
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

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.dynamicInfo?._id,
      context: LLM_SOCKET_CONTEXTS.dynamicInfos,
    });

    const dynamicInfoId = computed(() => config.value.dynamicInfo?._id);
    const title = computed(() => config.value.title || t('modals.createDynamicInfo.create.title'));
    const isDisabledIdField = computed(() => config.value.isDisabledIdField);

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.dynamicInfo,
      method: async () => {
        const result = await config.value.action?.(formToDynamicInfo(form.value));

        await config.value.afterSubmit?.(result);

        close();

        return result;
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
      submitting,
      chatShown,
      chatOptions,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>
