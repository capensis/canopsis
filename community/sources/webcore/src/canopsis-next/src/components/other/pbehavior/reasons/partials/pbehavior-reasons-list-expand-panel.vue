<template>
  <div class="planning-expand-panel">
    <v-tabs
      color="primary"
      class="secondary lighten-1"
      align-tabs="center"
      model-value="description"
    >
      <v-tab value="description">
        {{ t('common.description') }}
      </v-tab>
    </v-tabs>

    <v-window model-value="description">
      <c-window-expand-item value="description">
        <pre class="planning-expand-panel__pre">{{ pbehaviorReason.description || '-' }}</pre>
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    pbehaviorReason: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();

    return { t };
  },
};
</script>

<style lang="scss" scoped>
.planning-expand-panel {
  &__pre {
    margin: 0;
    white-space: pre-wrap;
  }
}
</style>
