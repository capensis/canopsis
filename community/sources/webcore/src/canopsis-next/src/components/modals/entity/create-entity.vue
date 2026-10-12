<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ config.title || $t('modals.createEntity.create.title') }}
      </template>
      <template #text="">
        <entity-form
          v-model="form"
          :prepare-state-setting-form="prepareStateSettingForm"
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
import { ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { entityToForm, formToEntity } from '@/helpers/entities/entity/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

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

    const form = ref(entityToForm(config.value.entity));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.entity,
      method: async () => {
        await config.value.action?.(formToEntity(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const prepareStateSettingForm = entity => ({
      ...formToEntity(entity),
      connector: config.value.entity.connector,
      _id: entity._id,
    });

    return {
      config,
      form,
      submitting,

      close,
      prepareStateSettingForm,
      submitLabel,
      submit,
    };
  },
};
</script>
