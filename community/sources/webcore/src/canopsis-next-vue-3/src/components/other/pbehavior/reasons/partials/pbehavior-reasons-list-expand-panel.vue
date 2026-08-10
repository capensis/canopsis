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
      <v-window-item value="description" class="planning-expand-panel__tab-item">
        <v-card>
          <v-card-text>
            <pre class="planning-expand-panel__pre">{{ pbehaviorReason.description || '-' }}
</pre>
          </v-card-text>
        </v-card>
      </v-window-item>
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

  &__tab-item {
    padding: 16px;
  }
}
</style>
