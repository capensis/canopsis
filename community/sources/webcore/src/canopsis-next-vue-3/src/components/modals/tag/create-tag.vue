<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <tag-form
          :form="values"
          :is-imported="isImported"
          :is-new="isNew"
          @update:form="updateForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn :disabled="isSubmitting" :loading="isSubmitting" color="primary" type="submit">
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { tagToForm, formToTag } from '@/helpers/entities/tag/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import TagForm from '@/components/other/tag/form/tag-form.vue';

export default {
  components: {
    TagForm,
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
    const config = computed(() => props.modal.config ?? {});
    const tag = computed(() => config.value.tag);
    const isNew = computed(() => !tag.value?._id);
    const title = computed(() => config.value.title || t('modals.createTag.create.title'));
    const isImported = computed(() => config.value.isImported);
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const form = useForm({
      initialValues: tagToForm(tag.value),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateForm = (value) => (values.value = value);

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToTag(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      values,
      title,
      isNew,
      isImported,
      isSubmitting,
      submitLabel,

      submit,
      updateForm,
      closeModal,
    };
  },
};
</script>
