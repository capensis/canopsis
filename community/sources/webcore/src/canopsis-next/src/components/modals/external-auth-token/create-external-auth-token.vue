<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <external-auth-token-form
          v-model="form"
          :rule-id="ruleId"
        />
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
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { externalAuthTokenToForm, formToExternalAuthToken } from '@/helpers/entities/external-auth-token/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ExternalAuthTokenForm from '@/components/other/external-auth-token/form/external-auth-token-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createExternalAuthToken,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ExternalAuthTokenForm,
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

    const form = ref(externalAuthTokenToForm(config.value.externalAuthToken));

    const title = computed(() => config.value.title || t('modals.createExternalAuthToken.create.title'));
    const ruleId = computed(() => config.value.externalAuthToken?._id);

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.externalAuthToken,
      method: async () => {
        const data = await config.value.action?.(formToExternalAuthToken(form.value));

        close();

        return data;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      config,

      form,

      ruleId,
      submitting,

      title,

      submitLabel,
      submit,
      close,
    };
  },
};
</script>
