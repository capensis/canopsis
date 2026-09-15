<template>
  <v-layout class="user-auth-key-field gap-2" align-center justify-space-between>
    <span class="user-auth-key-field__value">
      {{ displayedValue }}
    </span>
    <v-flex shrink>
      <c-action-btn
        v-if="value"
        :icon="visibilityIcon"
        :tooltip="visibilityTooltip"
        top
        @click="toggleVisibility"
      />
      <c-copy-btn
        v-if="value"
        :value="value"
        :tooltip="$t('common.copyToClipboard')"
        top
        @success="showCopySuccessPopup"
        @error="showCopyErrorPopup"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const { t } = useI18n();
    const popups = usePopups();

    const isShown = ref(false);

    const displayedValue = computed(() => {
      if (!props.value) {
        return '';
      }

      return isShown.value ? props.value : '•'.repeat(props.value.length);
    });

    const visibilityIcon = computed(() => (isShown.value ? 'visibility_off' : 'visibility'));
    const visibilityTooltip = computed(() => t(isShown.value ? 'common.hide' : 'common.show'));

    const toggleVisibility = () => isShown.value = !isShown.value;
    const showCopySuccessPopup = () => popups.success({ text: t('success.authKeyCopied') });
    const showCopyErrorPopup = () => popups.error({ text: t('errors.default') });

    return {
      displayedValue,
      visibilityIcon,
      visibilityTooltip,
      toggleVisibility,
      showCopySuccessPopup,
      showCopyErrorPopup,
    };
  },
};
</script>

<style lang="scss" scoped>
.user-auth-key-field {
  &__value {
    word-break: break-all;
  }
}
</style>
