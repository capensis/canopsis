<template>
  <modal-wrapper close>
    <template
      v-if="!config.hideTitle"
      #title=""
    >
      <span>{{ title }}</span>
    </template>
    <template
      v-if="sanitizedText || sanitizedAlertText"
      #text=""
    >
      <v-alert
        v-if="sanitizedAlertText"
        class="mb-2"
        type="warning"
      >
        <span v-html="sanitizedAlertText" />
      </v-alert>
      <span v-html="sanitizedText" class="text-subtitle-1 pre-wrap" />
    </template>
    <template #actions="">
      <v-layout
        wrap
        justify-center
      >
        <v-btn
          color="error"
          @click="cancel"
        >
          {{ $t('common.no') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="ml-2"
          color="primary"
          @click.prevent="submit"
        >
          {{ $t('common.yes') }}
        </v-btn>
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { ref, computed, onBeforeUnmount } from 'vue';

import { MODALS } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';

import { useInnerModal } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.confirmation,
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, modals } = useInnerModal(props);

    const submitted = ref(false);
    const cancelled = ref(false);
    const submitting = ref(false);

    const title = computed(() => (config.value.title ?? t('common.confirmation')));

    const sanitizedText = computed(() => (config.value.text ? sanitizeHtml(config.value.text) : ''));
    const sanitizedAlertText = computed(() => (config.value.alert ? sanitizeHtml(config.value.alert) : ''));

    const isDisabled = computed(() => submitting.value);

    const cancel = () => {
      cancelled.value = true;
      modals.hide();
    };

    const submit = async () => {
      if (config.value.action) {
        submitting.value = true;
        try {
          await config.value.action();
        } finally {
          submitting.value = false;
        }
      }

      submitted.value = true;
      modals.hide();
    };

    onBeforeUnmount(() => {
      if (!submitted.value && config.value.cancel) {
        config.value.cancel(cancelled.value);
      }
    });

    return {
      config,
      title,
      sanitizedText,
      sanitizedAlertText,
      submitting,
      isDisabled,
      cancel,
      submit,
    };
  },
};
</script>
