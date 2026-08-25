<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-payload-textarea-field
          :value="text"
          :required="required"
          :label="config.label"
          :name="config.name || 'text'"
          :variables="config.variables"
          autofocus
          @update:value="text = $event"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ $t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed, ref } from 'vue';

import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const modals = useModals();

const submitting = ref(false);
const text = ref(props.modal.config?.text ?? '');

const config = computed(() => props.modal.config ?? {});
const title = computed(() => config.value.title ?? t('modals.textEditor.title'));
const required = computed(() => (
  config.value.rules?.required ||
  config.value.rules?.includes?.('required') ||
  config.value.required
));
const isDisabled = computed(() => required.value && !text.value);

const close = () => {
  modals.hide({ id: props.modal.id });
};

const submit = async () => {
  if (isDisabled.value) {
    return;
  }

  submitting.value = true;

  try {
    await config.value.action?.(normalizeHandlebarsNbsp(text.value));

    close();
  } finally {
    submitting.value = false;
  }
};
</script>
