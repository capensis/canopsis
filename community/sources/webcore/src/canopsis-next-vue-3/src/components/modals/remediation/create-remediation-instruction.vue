<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        {{ title }}
      </template>
      <template #text>
        <remediation-instruction-approval-alert
          v-if="hasApproval && isChangesByCurrentUser"
          :user-name="alertUserName"
          :comment="alertComment"
          :dismissed="isChangesDismissed"
          class="mb-3"
        />
        <remediation-instruction-form
          v-model:form="form"
          :disabled="disabled"
          :is-new="isNew"
          :required-approve="requiredInstructionApprove"
          :rule-id="config.remediationInstruction?._id"
        />
        <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.attrs" />
      </template>
      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting || chatOptions.attrs.pending"
          :loading="submitting"
          color="primary"
          type="submit"
          @click="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, toRef } from 'vue';

import {
  LLM_SOCKET_CONTEXTS,
  MODALS,
  TEMPLATE_TESTING_TEST_TYPES,
} from '@/constants';

import {
  formToRemediationInstructionRequest,
  remediationInstructionErrorsToForm,
  remediationInstructionToFullForm,
} from '@/helpers/entities/remediation/instruction/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useAuth } from '@/hooks/auth';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useInfo } from '@/hooks/store/modules/info';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import RemediationInstructionForm from '@/components/other/remediation/instructions/form/remediation-instruction-form.vue';
import RemediationInstructionApprovalAlert from '@/components/other/remediation/instructions/partials/approval-alert.vue';

export default {
  name: MODALS.createRemediationInstruction,
  components: {
    AiChatSidebar,
    RemediationInstructionForm,
    RemediationInstructionApprovalAlert,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.instruction;

    const { t } = useI18n();
    const { currentUser } = useAuth();
    const { requiredInstructionApprove } = useInfo();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const close = () => modals.hide({ id: props.modal.id });

    const { values: form, handleSubmit, isSubmitting: submitting } = useForm({
      initialValues: remediationInstructionToFullForm(config.value.remediationInstruction),
      errorsToValidation: err => remediationInstructionErrorsToForm(err, form.value),
    });

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.remediationInstruction?._id,
      context: LLM_SOCKET_CONTEXTS.instruction,
    });

    const title = computed(
      () => config.value.title || t('modals.createRemediationInstruction.create.title'),
    );
    const disabled = computed(() => config.value.disabled);
    const isNew = computed(() => !config.value.remediationInstruction?._id);
    const approval = computed(() => config.value.remediationInstruction?.approval);
    const hasApproval = computed(() => !!approval.value);
    const isChangesDismissed = computed(() => !!approval.value?.dismissed_by);
    const isChangesByCurrentUser = computed(
      () => approval.value?.requested_by?._id === currentUser.value._id,
    );

    const alertUserName = computed(() => {
      const { dismissed_by: dismissedBy, requested_by: requestedBy } = approval.value ?? {};
      return dismissedBy?.display_name ?? requestedBy?.display_name;
    });

    const alertComment = computed(() => approval.value?.dismiss_comment ?? approval.value?.comment);

    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToRemediationInstructionRequest(form.value));

      await config.value.afterSubmit?.(result);

      close();

      return result;
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      disabled,
      approval,
      hasApproval,
      isChangesDismissed,
      isChangesByCurrentUser,
      alertUserName,
      alertComment,
      submitting,
      chatShown,
      chatOptions,
      requiredInstructionApprove,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>
