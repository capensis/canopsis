<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <dynamic-info-template-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { templateToForm, formToTemplate } from '@/helpers/entities/dynamic-info/template/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import DynamicInfoTemplateForm from '@/components/other/dynamic-info/form/dynamic-info-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDynamicInfoTemplate,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoTemplateForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(templateToForm(config.value.template ?? {}));

    const title = computed(() => (
      config.value.title || t('modals.createDynamicInfoTemplate.create.title')
    ));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.template,
      method: async () => {
        await config.value.action?.(formToTemplate(form.value));

        close();

        return form.value;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitting,
      submit,
      submitLabel,
      close,
    };
  },
};
</script>
