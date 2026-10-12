<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ t('modals.textEditorPaste.title') }}</span>
    </template>

    <template #text>
      <span class="text-subtitle-1">
        {{ t('modals.textEditorPaste.text') }}
      </span>
    </template>

    <template #actions>
      <v-btn variant="text" @click="cancel">
        {{ t('common.cancel') }}
      </v-btn>

      <v-btn color="primary" variant="outlined" @click="submit('text')">
        {{ t('modals.textEditorPaste.pasteAsPlainText') }}
      </v-btn>

      <v-btn color="primary" @click="submit('html')">
        {{ t('modals.textEditorPaste.pasteAsHtml') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, onBeforeUnmount, ref } from 'vue';

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

    const config = computed(() => props.modal.config ?? {});

    /**
     * Hides the paste confirmation modal.
     */
    const close = () => modals.hide({ id: props.modal.id });

    /**
     * Cancels paste insertion and closes the modal.
     *
     * Marks the modal as cancelled so `config.cancel` can restore the editor selection without inserting.
     */
    const cancel = () => {
      cancelled.value = true;
      close();
    };

    /**
     * Applies the chosen paste action and closes the modal.
     *
     * @param {'html' | 'text'} action - Insert as HTML or as plain text.
     */
    const submit = (action) => {
      submitted.value = true;
      config.value.action?.(action);
      close();
    };

    onBeforeUnmount(() => {
      if (!submitted.value && config.value.cancel) {
        config.value.cancel(cancelled.value);
      }
    });

    return {
      t,
      cancel,
      submit,
    };
  },
};
</script>
