<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createEntityPauseEvent.title') }}</span>
      </template>

      <template #text>
        <entity-pause-form :with-inherited="withInherited" />
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.create') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { entityPauseFormToRequest } from '@/helpers/entities/entity/pause';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityPauseForm from '@/components/forms/entity-pause-form.vue';

export default {
  name: MODALS.createEntityPauseEvent,

  components: { EntityPauseForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const form = useForm({
      initialValues: {
        comment: '',
        type: '',
        inherited: false,
        ...config.value.initialValues,
      },
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;
    const withInherited = computed(() => config.value.withInherited ?? true);

    const submit = handleSubmit(async () => {
      await config.value.action?.(entityPauseFormToRequest(values.value));
      close();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    return {
      t,
      values,
      withInherited,
      submitting,
      submit,
      close,
    };
  },
};
</script>
