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
            @cancel="chatCancelPending"
          />
          <patterns-form
            v-model="form"
            v-bind="patternsProps"
            autofocus
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
import { omit } from 'lodash';

import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY, LLM_SOCKET_CONTEXTS } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';
import PatternsForm from '@/components/forms/patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PatternsForm, ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(filterToForm(config.value.filter));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToFilter(form.value, PATTERNS_FIELDS, config.value.corporate));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));

    const patternsProps = computed(() => omit(config.value, ['title', 'action']));

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancelPending: chatCancelPending,
    } = useAiChatForm({
      form,
      modalId: props.modal.id,
      ruleId: props.modal.config?.filter?._id,
      context: LLM_SOCKET_CONTEXTS.widgetFilter,
    });

    return {
      title,
      form,
      patternsProps,
      isDisabled,
      submitting,
      close,
      submit,
      chatPending,
      chatPendingTexts,
      chatCancelPending,
    };
  },
};
</script>
