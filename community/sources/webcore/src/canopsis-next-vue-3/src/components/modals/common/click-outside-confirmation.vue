<template>
  <c-modal-wrapper>
    <template #title>
      <span>{{ t('modals.clickOutsideConfirmation.title') }}</span>
    </template>

    <template #text>
      <span class="text-subtitle-1">
        {{ t('modals.clickOutsideConfirmation.text') }}
      </span>
    </template>

    <template #actions>
      <v-btn
        variant="text"
        @click="close"
      >
        {{ t('modals.clickOutsideConfirmation.buttons.backToForm') }}
      </v-btn>

      <v-btn
        :loading="submitting"
        :disabled="submitting"
        color="primary"
        variant="outlined"
        @click.prevent="submit(false)"
      >
        {{ t('modals.clickOutsideConfirmation.buttons.dontSave') }}
      </v-btn>

      <v-btn
        :loading="submitting"
        :disabled="submitting"
        color="primary"
        @click.prevent="submit(true)"
      >
        {{ t('modals.clickOutsideConfirmation.buttons.save') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const submitting = ref(false);
    const config = computed(() => props.modal.config ?? {});

    const close = () => modals.hide({ id: props.modal.id });

    const submit = async (confirmed) => {
      submitting.value = true;

      try {
        await config.value.action?.(confirmed);

        close();
      } finally {
        submitting.value = false;
      }
    };

    return {
      t,
      submitting,
      close,
      submit,
    };
  },
};
</script>
