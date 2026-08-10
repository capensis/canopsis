<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <entity-info-property-form :is-new="isNew" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
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
import { computed } from 'vue';

import EntityInfoPropertyForm from '@/components/other/entity-info-property/form/entity-info-property-form.vue';
import {
  entityInfoPropertyToForm,
  formToEntityInfoProperty,
} from '@/helpers/entities/entity-info-property/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  components: {
    EntityInfoPropertyForm,
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
    const entityInfoProperty = computed(() => config.value.entityInfoProperty);
    const isNew = computed(() => !entityInfoProperty.value);
    const title = computed(
      () => config.value.title || t('modals.createEntityInfoProperty.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: entityInfoPropertyToForm(entityInfoProperty.value),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToEntityInfoProperty(values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      isNew,
      submitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
