<template>
  <c-tooltip location="left">
    <template #activator="{ props: activatorProps }">
      <v-btn
        color="secondary"
        size="x-large"
        icon
        v-bind="activatorProps"
        @click.stop="$emit('click', $event)"
      >
        <c-icon>refresh</c-icon>
      </v-btn>
    </template>
    <span>{{ t('common.refresh') }}</span>
  </c-tooltip>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['click'],

  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>
