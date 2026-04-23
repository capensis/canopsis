<template>
  <v-form class="position-relative" @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <c-progress-overlay :pending="defaultLlmPending" />
        <llm-form v-model="form" :is-new="isNew" :default-llm="defaultLlm" />
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
          :disabled="isDisabled"
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
import { computed, ref, onMounted } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { llmToForm, formToLlm } from '@/helpers/entities/llm/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

import LlmForm from '@/components/other/llm/form/llm-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createLlm,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    LlmForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const { fetchDefaultLlmWithoutStore } = useLlm();

    const form = ref(llmToForm(config.value.llm));
    const defaultLlm = ref(null);

    const isNew = computed(() => !config.value.llm?._id);

    const title = computed(() => (
      config.value.title || (
        isNew.value
          ? t('modals.createLlm.create.title')
          : t('modals.createLlm.edit.title')
      )
    ));

    const submitLabel = computed(() => (
      isNew.value ? t('common.add') : t('common.submit')
    ));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToLlm(form.value));

        close();
      },
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({ form, submit, close });

    const { pending: defaultLlmPending, handler: fetchDefaultLlm } = usePendingHandler(async () => {
      const { data: [fetchedDefaultLlm] = [] } = await fetchDefaultLlmWithoutStore() || {};

      if (!fetchedDefaultLlm || fetchedDefaultLlm._id === config.value.llm?._id) {
        return;
      }

      defaultLlm.value = fetchedDefaultLlm;

      if (!defaultLlm.value && isNew.value) {
        form.value.default = true;

        updateOriginalForm();
      }
    });

    onMounted(fetchDefaultLlm);

    return {
      form,

      isNew,
      title,
      submitLabel,

      isDisabled,
      submitting,

      submit,
      close,

      defaultLlmPending,
      defaultLlm,
    };
  },
};
</script>
