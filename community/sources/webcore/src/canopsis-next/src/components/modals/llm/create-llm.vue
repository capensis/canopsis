<template>
  <v-form class="position-relative" @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        {{ title }}
      </template>

      <template #text>
        <c-progress-overlay :pending="defaultLlmPending" />
        <llm-form
          :form="values"
          :is-new="isNew"
          :default-llm="defaultLlm"
          @update:form="updateForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { llmToForm, formToLlm } from '@/helpers/entities/llm/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLlm } from '@/hooks/store/modules/llm';

import LlmForm from '@/components/other/llm/form/llm-form.vue';

export default {
  components: {
    LlmForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchDefaultLlmWithoutStore } = useLlm();
    const config = computed(() => props.modal.config ?? {});
    const llm = computed(() => config.value.llm);
    const defaultLlm = ref(null);
    const isNew = computed(() => !llm.value?._id);

    const title = computed(
      () => config.value.title ||
        t(isNew.value ? 'modals.createLlm.create.title' : 'modals.createLlm.edit.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.add' : 'common.submit'));

    const form = useForm({
      initialValues: llmToForm(llm.value),
    });
    const { values, handleSubmit, isSubmitting, reset } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateForm = value => (values.value = value);

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToLlm(values.value));

      closeModal();
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    const { pending: defaultLlmPending, handler: fetchDefaultLlm } = usePendingHandler(async () => {
      const { data: [fetchedDefaultLlm] = [] } = (await fetchDefaultLlmWithoutStore()) || {};

      if (!fetchedDefaultLlm || fetchedDefaultLlm._id === llm.value?._id) {
        if (isNew.value && !fetchedDefaultLlm) {
          const nextValues = {
            ...values.value,
            default: true,
          };

          values.value = nextValues;
          reset(nextValues);
          updateOriginalForm();
        }

        return;
      }

      defaultLlm.value = fetchedDefaultLlm;
    });

    onMounted(fetchDefaultLlm);

    return {
      t,
      values,
      isNew,
      title,
      submitLabel,
      isSubmitting,
      defaultLlmPending,
      defaultLlm,

      submit,
      updateForm,
      closeModal,
    };
  },
};
</script>
