<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ config.title }}</span>
    </template>

    <template #text>
      <c-information-block :title="t('view.sharedViewUrl')">
        <c-text-field
          :model-value="config.url"
          class="share-view__url-field"
          readonly
        >
          <template #append-inner>
            <c-copy-btn
              :value="config.url"
              :tooltip="t('common.copyLink')"
              @success="showCopySuccessPopup"
              @error="showCopyErrorPopup"
            />
          </template>
        </c-text-field>
      </c-information-block>
    </template>

    <template #actions>
      <v-btn variant="text" @click="closeModal">
        {{ t('common.close') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed } from 'vue';

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

    const closeModal = () => modals.hide({ id: props.modal.id });
    const showCopySuccessPopup = () => popups.success({ text: t('success.linkCopied') });
    const showCopyErrorPopup = () => popups.error({ text: t('errors.default') });

    return {
      t,
      config,
      closeModal,
      showCopySuccessPopup,
      showCopyErrorPopup,
    };
  },
};
</script>

<style lang="scss" scoped>
.share-view__url-field {
  :deep(.v-field__append-inner) {
    align-items: center;
    padding-top: 0;
  }
}
</style>
