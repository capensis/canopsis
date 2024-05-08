<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <c-duration-field
          v-model="form"
          :label="label"
          :units="units"
          required
        />
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

import { durationToForm } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.duration,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(durationToForm(config.value.duration));

    const title = computed(() => config.value.title ?? t('common.duration'));
    const label = computed(() => config.value.label ?? t('common.duration'));
    const units = computed(() => (
      config.value.units?.map(unit => ({ ...unit, text: tc(unit.text, form.value.value) }))
    ));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form);
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      label,
      units,
      isDisabled,
      submitting,

      close,
      submit,
    };
  },
};
</script>
