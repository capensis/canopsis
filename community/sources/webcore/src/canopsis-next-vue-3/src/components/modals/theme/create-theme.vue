<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <theme-form />
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

import { formToTheme, themeToForm } from '@/helpers/entities/theme/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import ThemeForm from '@/components/other/theme/form/theme-form.vue';

export default {
  name: MODALS.createTheme,

  components: { ThemeForm },

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
    const form = useForm({
      initialValues: themeToForm(config.value.theme),
    });
    const { values, handleSubmit } = form;

    const title = computed(() => config.value.title ?? t('modals.createTheme.create.title'));
    const closeModal = () => modals.hide({ id: props.modal.id });

    const { submit, submitting, submitLabel } = useSubmittableForm({
      item: computed(() => config.value.theme),
      method: handleSubmit(async () => {
        const result = await config.value.action?.(formToTheme(values.value));

        await config.value.afterSubmit?.(result);

        closeModal();

        return result;
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
      submitting,
      closeModal,
      submitLabel,
      submit,
    };
  },
};
</script>
