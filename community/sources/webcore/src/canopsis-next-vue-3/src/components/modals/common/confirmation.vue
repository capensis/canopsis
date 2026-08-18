<template>
  <c-modal-wrapper close>
    <template v-if="!config.hideTitle" #title>
      <span>{{ title }}</span>
    </template>

    <template v-if="sanitizedText || sanitizedAlertText" #text>
      <v-alert v-if="sanitizedAlertText" class="mb-2 pre-line" type="warning">
        <span v-html="sanitizedAlertText" />
      </v-alert>

      <span v-html="sanitizedText" class="text-subtitle-1 pre-wrap" />
    </template>

    <template #actions>
      <c-layout class="flex-grow-1" justify="center" gap="2">
        <v-btn
          :variant="config.cancelOutlined ? 'outlined' : 'elevated'"
          color="error"
          @click="cancel"
        >
          {{ config.cancelText || t('common.no') }}
        </v-btn>

        <v-btn
          :loading="submitting"
          :disabled="submitting"
          variant="elevated"
          color="primary"
          @click.prevent="submit"
        >
          {{ config.actionText || t('common.yes') }}
        </v-btn>
      </c-layout>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, onBeforeUnmount, ref } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

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

    const submitted = ref(false);
    const cancelled = ref(false);
    const submitting = ref(false);

    const config = computed(() => props.modal.config ?? {});
    const title = computed(() => config.value.title ?? t('common.confirmation'));
    const sanitizedText = computed(() =>
      config.value.text ? sanitizeHtml(config.value.text) : '',
    );
    const sanitizedAlertText = computed(() =>
      config.value.alert ? sanitizeHtml(config.value.alert) : '',
    );

    const close = () => modals.hide({ id: props.modal.id });

    const cancel = () => {
      cancelled.value = true;
      close();
    };

    const submit = async () => {
      try {
        submitting.value = true;

        await config.value.action?.();

        submitted.value = true;
        close();
      } finally {
        submitting.value = false;
      }
    };

    onBeforeUnmount(() => {
      if (!submitted.value && config.value.cancel) {
        config.value.cancel(cancelled.value);
      }
    });

    return {
      config,
      submitting,
      title,
      sanitizedText,
      sanitizedAlertText,
      cancel,
      submit,
      t,
    };
  },
};
</script>
