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
          <kpi-filter-form v-model="form" />
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

import { LLM_SOCKET_CONTEXTS, MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import KpiFilterForm from '@/components/other/kpi/filters/form/kpi-filter-form.vue';
import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createKpiFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { KpiFilterForm, ModalWrapper, PatternProgress },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const { filter = {} } = props.modal.config;

    const form = ref({
      name: filter.name ?? '',
      patterns: filterPatternsToForm(filter, [PATTERNS_FIELDS.entity]),
    });

    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));

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
      ruleId: props.modal.config?.filter?._id,
      context: LLM_SOCKET_CONTEXTS.kpiFilter,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.({
          name: form.value.name,
          ...formFilterToPatterns(form.value.patterns, [PATTERNS_FIELDS.entity]),
        });
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
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
