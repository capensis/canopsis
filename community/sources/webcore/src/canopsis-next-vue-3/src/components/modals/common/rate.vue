<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ config.title }}</span>
      </template>

      <template #text>
        <c-layout justify="center">
          <span class="text-subtitle-1">
            {{ config.text }}
          </span>
        </c-layout>
        <rate-form v-model:form="form" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn :loading="submitting" color="primary" type="submit">
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import RateForm from '@/components/forms/rate.vue';
import { MODALS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';

export default {
  name: MODALS.rate,

  components: { RateForm },

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
    const form = ref({
      comment: '',
      rating: 5,
    });
    const closeModal = () => modals.hide({ id: props.modal.id });
    const { pending: submitting, handler: submit } = usePendingHandler(async () => {
      await config.value.action?.(form.value);

      closeModal();
    });

    return {
      t,
      config,
      form,
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>
