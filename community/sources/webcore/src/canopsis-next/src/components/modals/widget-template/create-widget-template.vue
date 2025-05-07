<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <widget-template-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
          color="primary"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { widgetTemplateToForm, formToWidgetTemplate } from '@/helpers/entities/widget/template/form';

import { useInnerModal } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import WidgetTemplateForm from '@/components/other/widget-template/form/widget-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWidgetTemplate,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    WidgetTemplateForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(widgetTemplateToForm(config.value.widgetTemplate));

    const title = computed(() => config.value.title ?? t('modals.createWidgetTemplate.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToWidgetTemplate(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submit,
      isDisabled,
      submitting,
    };
  },
};
</script>
