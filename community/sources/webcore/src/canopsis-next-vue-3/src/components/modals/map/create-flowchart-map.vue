<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <flowchart-map-form v-model:form="form" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn :disabled="isSubmitting" :loading="isSubmitting" color="primary" type="submit">
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { formToMap, mapToForm } from '@/helpers/entities/map/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import FlowchartMapForm from '@/components/other/map/form/flowchart-map-form.vue';

export default {
  name: MODALS.createFlowchartMap,

  components: { FlowchartMapForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const title = computed(() => config.value.title ?? t('modals.createFlowchartMap.create.title'));

    const {
      values: form,
      handleSubmit,
      isSubmitting,
    } = useForm({
      initialValues: mapToForm(config.value.map),
    });

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToMap(form.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close: closeModal,
    });

    return {
      t,
      form,
      title,
      isSubmitting,
      submit,
      closeModal,
    };
  },
};
</script>
