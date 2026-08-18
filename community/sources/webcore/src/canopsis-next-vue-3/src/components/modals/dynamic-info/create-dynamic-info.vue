<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <dynamic-info-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ $t('common.cancel') }}
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

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

export default {
  components: { DynamicInfoForm },

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
    const isNew = computed(() => !config.value.dynamicInfo);
    const title = computed(() => config.value.title || t('modals.createDynamicInfo.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const { values, handleSubmit, submitting } = useForm({
      initialValues: dynamicInfoToForm(config.value.dynamicInfo),
    });

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToDynamicInfo(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();
    });

    useFormConfirmableCloseModal({
      submit,

      form: values,
      close: closeModal,
    });

    return {
      config,
      submitting,
      isNew,
      title,
      submitLabel,
      values,
      submit,
      closeModal,
    };
  },
};
</script>
