<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <v-text-field
          v-model="text"
          v-validate="field.validationRules"
          v-bind="fieldProps"
          :error-messages="errors.collect(field.name)"
          autofocus
        />
        <c-alert
          v-if="config.alert"
          type="warning"
        >
          {{ config.alert.text }}
        </c-alert>
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
import { omit } from 'lodash';
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textFieldEditor,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { modals, config, close } = useInnerModal(props);

    const text = ref(config.value.field?.value ?? '');

    const form = ref({ text: text.value });

    const field = computed(() => config.value.field ?? { name: 'text', label: 'Text' });
    const fieldProps = computed(() => omit(field.value, ['validationRules', 'value']));

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        await config.value?.action?.(text.value);

        modals.hide();
      },
    });

    useFormConfirmableCloseModal({
      form: text,
      submit,
      close,
    });

    return {
      text,
      field,
      fieldProps,
      submit,
      submitting,
      config,
      isDisabled,
    };
  },
};
</script>
