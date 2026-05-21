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
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { widgetTemplateToForm, formToWidgetTemplate } from '@/helpers/entities/widget/template/form';
import { getWidgetTemplateModalTitle } from '@/helpers/entities/widget/template/modal-title';

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
    const { t, te } = useI18n();

    const form = ref(widgetTemplateToForm(config.value.widgetTemplate));

    const title = computed(() => {
      if (config.value.title) {
        return config.value.title;
      }

      return getWidgetTemplateModalTitle({
        t,
        te,
        type: form.value.type ?? config.value.widgetTemplate?.type,
        isEdit: !!config.value.widgetTemplate?._id,
      });
    });

    const { submit, isDisabled, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.widgetTemplate,
      method: async () => {
        await config.value.action?.(formToWidgetTemplate(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitLabel,
      submit,
      isDisabled,
      submitting,
    };
  },
};
</script>
