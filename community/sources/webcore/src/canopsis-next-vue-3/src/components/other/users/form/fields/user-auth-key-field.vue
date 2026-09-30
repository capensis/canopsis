<template>
  <c-layout class="user-auth-key-field" align="center" justify="space-between">
    <span class="user-auth-key-field__value">
      {{ displayedValue }}
    </span>

    <c-layout v-if="value" align="center" gap="1">
      <c-tooltip location="top">
        <template #activator="{ props: activatorProps }">
          <v-btn
            v-bind="activatorProps"
            size="small"
            variant="text"
            icon
            @click="toggleVisibility"
          >
            <c-icon>{{ visibilityIcon }}</c-icon>
          </v-btn>
        </template>

        <span>{{ visibilityTooltip }}</span>
      </c-tooltip>

      <c-tooltip location="top">
        <template #activator="{ props: activatorProps }">
          <v-btn
            v-bind="activatorProps"
            size="small"
            variant="text"
            icon
            @click="copy"
          >
            <c-icon>content_copy</c-icon>
          </v-btn>
        </template>

        <span>{{ t('common.copyToClipboard') }}</span>
      </c-tooltip>
    </c-layout>
  </c-layout>
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

    const toggleVisibility = () => (isShown.value = !isShown.value);

    const copy = async () => {
      try {
        await navigator.clipboard.writeText(props.value);
        popups.success({ text: t('success.authKeyCopied') });
      } catch (err) {
        console.error(err);
        popups.error({ text: t('errors.default') });
      }
    };

    return {
      t,
      displayedValue,
      visibilityIcon,
      visibilityTooltip,
      toggleVisibility,
      copy,
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
