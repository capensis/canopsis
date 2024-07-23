<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.launchEventsRecord.title') }}</span>
      </template>
      <template #text="">
        <c-alarm-patterns-field
          v-model="form"
          :attributes="attributes"
          name="patterns"
          required
          @input="errors.remove('patterns')"
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
          {{ $t('common.launch') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { ENTITY_PATTERN_FIELDS, MODALS, VALIDATION_DELAY } from '@/constants';

import { patternToForm } from '@/helpers/entities/pattern/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.startEventsRecord,
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
    const form = ref(patternToForm());
    const attributes = computed(() => [
      {
        value: ENTITY_PATTERN_FIELDS.lastEventDate,
        options: { disabled: true },
      },
    ]); // TODO: finish it

    const { config, close } = useInnerModal(props);
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
      attributes,
      isDisabled,
      submitting,

      submit,
      close,
    };
  },
};
</script>
