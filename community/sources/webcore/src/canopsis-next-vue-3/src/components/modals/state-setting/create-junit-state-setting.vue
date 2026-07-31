<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <junit-state-setting-form />
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
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  formToJunitStateSetting,
  junitStateSettingToForm,
} from '@/helpers/entities/state-setting/junit/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import JunitStateSettingForm from '@/components/other/state-setting/junit/form/junit-state-setting-form.vue';

export default {
  components: { JunitStateSettingForm },

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
    const title = computed(
      () => config.value.title || t('modals.createJunitStateSetting.title'),
    );

    const form = useForm({
      initialValues: junitStateSettingToForm(config.value.stateSetting),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToJunitStateSetting(values.value));

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
      isSubmitting,
      submit,
      closeModal,
    };
  },
};
</script>
