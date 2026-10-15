<template>
  <div class="planning-expand-panel">
    <v-tabs
      v-model="activeTab"
      color="primary"
      class="secondary lighten-1"
      align-tabs="center"
    >
      <v-tab :value="TABS.description">
        {{ t('common.description') }}
      </v-tab>
      <v-tab :value="TABS.type">
        {{ t('common.type') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <c-window-expand-item :value="TABS.description">
        <pre class="planning-expand-panel__pre">{{ pbehaviorType.description || '-' }}</pre>
      </c-window-expand-item>

      <c-window-expand-item :value="TABS.type">
        <c-layout gap="3" align="center">
          <v-chip
            v-if="pbehaviorType.icon_name"
            :color="pbehaviorType.color"
            class="planning-expand-panel__type-chip"
          >
            <v-icon :color="iconColor" :icon="pbehaviorType.icon_name" size="18" />
          </v-chip>
          <span>{{ typeLabel }}</span>
        </c-layout>
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';

import { useI18n } from '@/hooks/i18n';

const TABS = {
  description: 'description',
  type: 'type',
};

export default {
  props: {
    pbehaviorType: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const activeTab = ref(TABS.description);
    const iconColor = computed(() => getMostReadableTextColor(props.pbehaviorType.color, { level: 'AA', size: 'large' }),
    );
    const typeLabel = computed(() => props.pbehaviorType.type ? t(`pbehavior.types.types.${props.pbehaviorType.type}`) : '-',
    );

    return {
      TABS,
      activeTab,
      iconColor,
      typeLabel,
      t,
    };
  },
};
</script>

<style lang="scss" scoped>
.planning-expand-panel {
  &__pre {
    margin: 0;
    white-space: pre-wrap;
  }
  &__type-chip {
    width: 42px;
    height: 24px;
  }
}
</style>
