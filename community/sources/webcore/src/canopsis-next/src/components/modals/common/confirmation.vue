<template>
  <modal-wrapper close>
    <template
      v-if="!config.hideTitle"
      #title=""
    >
      <span>{{ title }}</span>
    </template>
    <template
      v-if="config.text"
      #text=""
    >
      <span class="text-subtitle-1 pre-wrap">{{ config.text }}</span>
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
          {{ config.cancelText || $t('common.no') }}
        </v-btn>
        <v-btn
          v-if="config.secondAction"
          :loading="submittingSecondAction"
          :disabled="isDisabledSecondAction"
          class="ml-2"
          color="primary"
          outlined
          @click.prevent="submitSecondAction"
        >
          {{ config.secondActionText }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="ml-2"
          color="primary"
          @click.prevent="submit"
        >
          {{ config.actionText || $t('common.yes') }}
        </v-btn>
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { ref, computed, onBeforeUnmount } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

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
    // Reactive data
    const submitted = ref(false);
    const cancelled = ref(false);

    // Composables
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    // Computed properties
    const title = computed(() => config.value.title ?? t('common.confirmation'));

    // Methods
    const cancel = () => {
      cancelled.value = true;
      close();
    };

    const { submit, submitting, isDisabled } = useSubmittableForm({
      method: async () => {
        await config.value.action?.();

        submitted.value = true;
        close();
      },
    });

    const {
      submit: submitSecondAction,
      submitting: submittingSecondAction,
      isDisabled: isDisabledSecondAction,
    } = useSubmittableForm({
      method: async () => {
        await config.value.secondAction?.();

        submitted.value = true;
        close();
      },
    });

    // Lifecycle
    onBeforeUnmount(() => {
      if (!submitted.value) {
        config.value.cancel?.(cancelled.value);
      }
    });

    return {
      config,
      submitting,
      isDisabled,
      title,
      cancel,
      submit,
      submitSecondAction,
      submittingSecondAction,
      isDisabledSecondAction,
    };
  },
};
</script>
