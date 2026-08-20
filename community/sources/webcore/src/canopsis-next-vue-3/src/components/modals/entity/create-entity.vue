<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        {{ title }}
      </template>

      <template #text>
        <entity-form
          :form="values"
          :prepare-state-setting-form="prepareStateSettingForm"
          @update:form="updateForm"
        />
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="closeModal">
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

import { MODALS } from '@/constants';

import { entityToForm, formToEntity } from '@/helpers/entities/entity/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

export default {
  name: MODALS.createEntity,

  components: {
    EntityForm,
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
    const entity = computed(() => config.value.entity);
    const title = computed(() => config.value.title || t('modals.createEntity.create.title'));
    const form = useForm({
      initialValues: entityToForm(entity.value),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateForm = value => (values.value = value);
    const submitLabel = computed(() => t(entity.value?._id ? 'common.save' : 'common.create'));

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToEntity(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    const prepareStateSettingForm = preparedEntity => ({
      ...formToEntity(preparedEntity),
      connector: entity.value?.connector,
      _id: preparedEntity._id,
    });

    return {
      t,
      title,
      values,
      submitting,
      submitLabel,
      closeModal,
      prepareStateSettingForm,
      updateForm,
      submit,
    };
  },
};
</script>
