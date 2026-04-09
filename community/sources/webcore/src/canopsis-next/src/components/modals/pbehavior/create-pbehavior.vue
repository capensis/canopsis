<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <div class="position-relative">
          <pattern-progress
            v-if="chatPending"
            :in-progress-text="chatPendingTexts.inProgress"
            :cancel-button-label="chatPendingTexts.cancel"
            @cancel="chatCancel"
          />
          <pbehavior-form
            v-model="form"
            :no-pattern="noPattern"
            :with-inherited="withInherited"
            :pbehavior-id="pbehaviorId"
            pbehavior-counter-type
          />
        </div>
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

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorToForm, formToPbehavior, pbehaviorToRequest } from '@/helpers/entities/pbehavior/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';
import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PbehaviorForm, ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const { pbehavior, timezone } = props.modal.config;

    const form = ref(pbehaviorToForm(pbehavior, null, timezone));

    const title = computed(() => config.value.title || t('modals.createPbehavior.create.title'));
    const noPattern = computed(() => !!config.value.noPattern);
    const withInherited = computed(() => !!config.value.withInherited);
    const pbehaviorId = computed(() => config.value.pbehavior?._id);

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancel: chatCancel,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: props.modal.id,
      ruleId: props.modal.config?.pbehavior?._id,
      context: LLM_SOCKET_CONTEXTS.pbehavior,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(
          pbehaviorToRequest(formToPbehavior(form.value, config.value.timezone)),
        );
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      noPattern,
      withInherited,
      pbehaviorId,
      isDisabled,
      submitting,
      chatPending,
      chatPendingTexts,
      chatCancel,
      submit,
      close,
    };
  },
};
</script>
