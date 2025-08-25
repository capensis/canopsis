<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <entity-info-property-form v-model="form" :is-new="isNew" />
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
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { entityInfoPropertyToForm } from '@/helpers/entities/entity-info-property/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import EntityInfoPropertyForm from '@/components/other/entity-info-property/form/entity-info-property-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntityInfoProperty,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EntityInfoPropertyForm,
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

    const form = ref(entityInfoPropertyToForm(config.value.entityInfoProperty));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title || t('modals.createEntityInfoProperty.create.title'));
    const isNew = computed(() => !config.value.entityInfoProperty);

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
