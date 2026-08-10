<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.addDynamicInfoInfosFromTemplate.title') }}</span>
      </template>

      <template #text>
        <c-select
          v-model="templateField"
          :disabled="templatesPending"
          :items="templates"
          :label="t('modals.addDynamicInfoInfosFromTemplate.fields.template')"
          :loading="templatesPending"
          item-title="title"
          item-value="_id"
          name="template"
          return-object
          required
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="templatesPending || submitting"
          :loading="templatesPending || submitting"
          color="primary"
          type="submit"
        >
          {{ t('modals.addDynamicInfoInfosFromTemplate.actions.addInfos') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useDynamicInfoTemplates } from '@/hooks/store/modules/dynamic-info-templates';

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchDynamicInfoTemplatesList } = useDynamicInfoTemplates();

    const templates = ref([]);

    const config = computed(() => props.modal.config ?? {});
    const form = useForm({
      initialValues: {
        template: null,
      },
    });
    const templateField = form.model('template');

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = form.handleSubmit(async () => {
      await config.value.action?.(form.values.value.template);

      closeModal();
    });

    useFormConfirmableCloseModal({
      submit,

      form: form.values,
      close: closeModal,
    });

    const { pending: templatesPending, handler: fetchTemplatesList } = usePendingHandler(
      async () => {
        templates.value = await fetchDynamicInfoTemplatesList();
      },
      true,
    );

    onMounted(fetchTemplatesList);

    return {
      t,
      templates,
      templatesPending,
      submitting: form.submitting,
      templateField,
      submit,
      closeModal,
    };
  },
};
</script>
