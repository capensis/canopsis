<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ $t('modals.infoPopupSetting.addInfoPopup.title') }}</span>
      </template>

      <template #text>
        <info-popup-form
          v-model:form="form"
          :columns="config.columns"
          :variables="config.variables"
        />
      </template>

      <template #actions>
        <v-btn :disabled="isSubmitting" variant="text" @click="closeModal">
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="isSubmitting"
          :disabled="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useModals } from '@/hooks/modals';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';

export default {
  name: MODALS.addInfoPopup,

  components: { InfoPopupForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const popup = computed(() => config.value.popup ?? {});
    const firstColumn = computed(() => config.value.columns?.[0]);
    const formContext = useForm({
      initialValues: {
        form: {
          column: popup.value.column ?? firstColumn.value?.value ?? firstColumn.value?.column,
          template: popup.value.template ?? '',
        },
      },
    });
    const { handleSubmit, isSubmitting, model, values } = formContext;
    const form = model('form');
    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(values.value.form);

      closeModal();
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close: closeModal,
    });

    return {
      config,
      form,
      isSubmitting,
      submit,
      closeModal,
    };
  },
};
</script>
