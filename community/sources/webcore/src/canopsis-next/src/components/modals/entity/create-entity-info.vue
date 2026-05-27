<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <entity-info-form
          v-model="form"
          :entity-info="config.entityInfo"
          :infos="config.infos"
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

import { entityInfoToForm, formToEntityInfo } from '@/helpers/entities/entity-info/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import EntityInfoForm from '@/components/other/entity/form/entity-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntityInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { EntityInfoForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(entityInfoToForm(config.value.entityInfo));

    const title = computed(() => config.value.title ?? t('modals.createEntityInfo.create.title'));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.entityInfo,
      isNewCheck: value => !value?.name,
      method: async () => {
        await config.value.action?.(formToEntityInfo(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      title,
      submit,
      submitting,
      submitLabel,
      close,
    };
  },
};
</script>
