<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createDynamicInfoInformation.create.title') }}</span>
      </template>
      <template #text="">
        <dynamic-info-information-form
          v-model="form"
          :existing-names="existingNames"
          :initial-name="initialName"
          :variables="variables"
          :copy-variables="copyVariables"
        />
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
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { dynamicInfoInformationToForm } from '@/helpers/entities/dynamic-info/information/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import DynamicInfoInformationForm from '@/components/other/dynamic-info/form/fields/dynamic-info-information-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoInformation,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    DynamicInfoInformationForm,
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

    const { info = {} } = config.value;

    const form = ref(dynamicInfoInformationToForm(info));

    const initialName = computed(() => config.value.info && config.value.info.name);
    const existingNames = computed(() => config.value.existingNames);
    const variables = computed(() => config.value.variables);
    const copyVariables = computed(() => config.value.copyVariables);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const data = await config.value.action?.(form.value);

        close();

        return data;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,

      existingNames,
      initialName,
      variables,
      copyVariables,
      isDisabled,
      submitting,
      submit,
    };
  },
};
</script>
