<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ config.title }}
      </template>
      <template #text="">
        <template-testing-test-form v-model="form" :is-new="isNew" />
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

import { templateTestingTestToForm } from '@/helpers/entities/template-testing-test/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import TemplateTestingTestForm from '@/components/other/template-testing/form/template-testing-test-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTemplateTestingTest,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    TemplateTestingTestForm,
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

    const form = ref(templateTestingTestToForm(config.value.templateTestingTest));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const isNew = computed(() => !config.value.templateTestingTest);

    return {
      form,

      isDisabled,
      submitting,

      isNew,

      submit,
      close,
    };
  },
};
</script>
