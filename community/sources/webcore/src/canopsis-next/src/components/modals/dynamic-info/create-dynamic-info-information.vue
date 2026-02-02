<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createDynamicInfoInformation.create.title') }}</span>
      </template>
      <template #text="">
        <div>
          <v-text-field
            v-model="form.name"
            v-validate="nameRules"
            :label="$t('common.name')"
            :error-messages="errors.collect('name')"
            name="name"
            autofocus
          />
          <c-mixed-field
            v-model="form.value"
            :label="$t('common.value')"
            :types="fieldTypes"
            name="value"
            required
          />
        </div>
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

import { MODALS, VALIDATION_DELAY, PATTERN_FIELD_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

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
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const fieldTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ];

    const form = ref({
      name: config.value?.info?.name ?? '',
      value: config.value?.info?.value ?? '',
    });

    const initialName = computed(() => config.value?.info?.name);
    const existingNames = computed(() => config.value.existingNames);

    const nameRules = computed(() => ({
      required: true,
      unique: {
        values: existingNames.value,
        initialValue: initialName.value,
      },
    }));

    const { submit, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(form.value);
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      nameRules,
      fieldTypes,
      isDisabled,
      submit,
      close,
      t,
    };
  },
};
</script>
