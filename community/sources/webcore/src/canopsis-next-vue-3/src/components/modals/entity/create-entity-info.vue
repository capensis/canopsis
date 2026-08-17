<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <entity-info-form
          :form="values"
          :entity-info="entityInfo"
          :infos="infos"
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

import { entityInfoToForm, formToEntityInfo } from '@/helpers/entities/entity-info/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import EntityInfoForm from '@/components/other/entity/form/entity-info-form.vue';

export default {
  name: MODALS.createEntityInfo,

  components: {
    EntityInfoForm,
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
    const entityInfo = computed(() => config.value.entityInfo);
    const infos = computed(() => config.value.infos ?? []);
    const title = computed(() => config.value.title ?? t('modals.createEntityInfo.create.title'));
    const form = useForm({
      initialValues: entityInfoToForm(entityInfo.value),
    });
    const { values, handleSubmit } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateForm = value => (values.value = value);

    const { submit, submitting, submitLabel } = useSubmittableForm({
      item: entityInfo,
      isNewCheck: value => !value?.name,
      method: handleSubmit(async () => {
        await config.value.action?.(formToEntityInfo(values.value));

        closeModal();
      }),
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      values,
      entityInfo,
      infos,
      submitting,
      submitLabel,
      closeModal,
      updateForm,
      submit,
    };
  },
};
</script>
