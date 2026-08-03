<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ config.title }}</span>
      </template>

      <template #text>
        <v-alert class="mb-3 pre-line" type="warning">
          <span v-html="sanitizedText" class="text-body-2" />
        </v-alert>

        <div class="my-3">
          <p v-html="sanitizedPhraseText" class="mb-2" />
          <code
            v-text="originalPhrase"
            class="text-black bg-grey-lighten-2 d-inline-block pa-1"
          />
        </div>

        <c-text-field
          v-model="phrase"
          :label="t('modals.confirmationPhrase.phrase')"
          class="mt-2"
          autofocus
        />

        <component
          v-if="config.additionalForm?.component"
          :is="config.additionalForm.component"
          v-model="additionalForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :loading="submitting"
          :disabled="submitting || !phrasesEqual"
          color="primary"
          type="submit"
        >
          {{ t('common.yes') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

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
    const phrase = ref('');
    const submitting = ref(false);
    const config = computed(() => props.modal.config ?? {});
    const additionalForm = ref(config.value.additionalForm?.form);
    const originalPhrase = computed(() => config.value.phrase);
    const phrasesEqual = computed(() => phrase.value === originalPhrase.value);
    const sanitizedText = computed(() => sanitizeHtml(config.value.text || ''));
    const sanitizedPhraseText = computed(() => sanitizeHtml(config.value.phraseText || ''));

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = async () => {
      if (!phrasesEqual.value) {
        return;
      }

      try {
        submitting.value = true;

        await config.value.action?.(additionalForm.value);

        closeModal();
      } finally {
        submitting.value = false;
      }
    };

    return {
      t,
      phrase,
      additionalForm,
      config,
      originalPhrase,
      phrasesEqual,
      sanitizedText,
      sanitizedPhraseText,
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>
