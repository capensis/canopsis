<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ config.title }}</span>
    </template>

    <template #text>
      <v-alert v-if="config.warningText" class="mb-4" type="warning">
        {{ config.warningText }}
      </v-alert>

      <c-text-field
        :model-value="config.value"
        :label="config.label"
        class="copy-value-modal__field"
        readonly
      >
        <template #append-inner>
          <c-copy-btn
            :value="config.value"
            :tooltip="copyTooltip"
            @success="showCopySuccessPopup"
            @error="showCopyErrorPopup"
          />
        </template>
      </c-text-field>
    </template>

    <template #actions>
      <v-btn variant="text" @click="closeModal">
        {{ t('common.close') }}
      </v-btn>

      <v-btn color="primary" @click="copyValue">
        {{ t('common.copy') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { writeTextToClipboard } from '@/helpers/clipboard';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';

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
    const popups = usePopups();
    const config = computed(() => props.modal.config ?? {});
    const copyTooltip = computed(
      () => config.value.copyTooltip ?? t('common.copyToClipboard'),
    );

    const closeModal = () => modals.hide({ id: props.modal.id });
    const showCopySuccessPopup = () => popups.success({
      text: config.value.copySuccessText ?? t('success.default'),
    });
    const showCopyErrorPopup = () => popups.error({ text: t('errors.default') });
    const copyValue = async () => {
      try {
        await writeTextToClipboard(config.value.value);
        showCopySuccessPopup();
      } catch {
        showCopyErrorPopup();
      }
    };

    return {
      t,
      config,
      copyTooltip,
      closeModal,
      copyValue,
      showCopySuccessPopup,
      showCopyErrorPopup,
    };
  },
};
</script>

<style lang="scss" scoped>
.copy-value-modal__field {
  :deep(.v-field__append-inner) {
    align-items: center;
    padding-top: 0;
  }
}
</style>
