<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <entity-form v-model="form" :prepare-state-setting-form="prepareStateSettingForm" />
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

import { entityToForm, formToEntity } from '@/helpers/entities/entity/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EntityForm,
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
    const { t } = useI18n();

    const form = ref(entityToForm(config.value.entity));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToEntity(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title || t('modals.createEntity.create.title'));

    const prepareStateSettingForm = entity => ({
      ...formToEntity(entity),
      connector: config.value.entity.connector,
      _id: entity._id,
    });

    return {
      title,
      config,
      form,
      isDisabled,
      submitting,

      close,
      prepareStateSettingForm,
      submit,
    };
  },
};
</script>
