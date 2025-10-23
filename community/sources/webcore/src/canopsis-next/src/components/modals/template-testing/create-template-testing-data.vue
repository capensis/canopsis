<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <template-testing-data-form v-model="form" :is-new="isNew" />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.save') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { templateTestingDataToForm, formToTemplateTestingData } from '@/helpers/entities/template-testing-data/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import TemplateTestingDataForm from '@/components/other/template-testing/form/template-testing-data-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTemplateTestingData,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    TemplateTestingDataForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(templateTestingDataToForm(config.value.templateTestingData));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToTemplateTestingData(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title || t('modals.createTemplateTestingData.create.title'));

    const isNew = computed(() => !config.value.templateTestingData?._id);

    return {
      form,

      isDisabled,
      submitting,

      title,
      isNew,

      submit,
      close,
    };
  },
};
</script>
