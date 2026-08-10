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
      <v-tab :value="TABS.periods">
        {{ t('common.periods') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item :value="TABS.description" class="planning-expand-panel__tab-item">
        <v-card>
          <v-card-text>
            <pre class="planning-expand-panel__pre">{{
              pbehaviorException.description || '-'
            }}</pre>
          </v-card-text>
        </v-card>
      </v-window-item>

      <v-window-item :value="TABS.periods" class="planning-expand-panel__tab-item">
        <v-card>
          <v-card-text>
            <c-advanced-data-table
              :items="pbehaviorException.exdates || []"
              :headers="headers"
              item-key="begin"
              hide-actions
              no-pagination
            >
              <template #begin="{ item }">
                {{ formatDate(item.begin) }}
              </template>
              <template #end="{ item }">
                {{ formatDate(item.end) }}
              </template>
              <template #type="{ item }">
                <c-layout v-if="item.type" gap="2" align="center">
                  <v-chip
                    v-if="item.type.icon_name"
                    :color="item.type.color"
                    class="planning-expand-panel__type-chip"
                  >
                    <v-icon
                      :color="getIconColor(item.type.color)"
                      :icon="item.type.icon_name"
                      size="18"
                    />
                  </v-chip>
                  <span>{{ item.type.name || item.type }}</span>
                </c-layout>
                <span v-else>-</span>
              </template>
            </c-advanced-data-table>
          </v-card-text>
        </v-card>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { DATETIME_FORMATS } from '@/constants';
import { getMostReadableTextColor } from '@/helpers/color';
import { convertDateToString } from '@/helpers/date/date';
import { useI18n } from '@/hooks/i18n';

const TABS = {
  description: 'description',
  periods: 'periods',
};

export default {
  props: {
    pbehaviorException: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();

    const activeTab = ref(TABS.description);
    const headers = computed(() => [
      { value: 'begin', text: t('common.begin'), sortable: false },
      { value: 'end', text: t('common.end'), sortable: false },
      { value: 'type', text: t('common.type'), sortable: false },
    ]);
    const formatDate = value => convertDateToString(value, DATETIME_FORMATS.long);
    const getIconColor = color => getMostReadableTextColor(color, { level: 'AA', size: 'large' });

    return {
      TABS,
      activeTab,
      headers,
      formatDate,
      getIconColor,
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

  &__tab-item {
    padding: 16px;
  }

  &__type-chip {
    width: 42px;
    height: 24px;
  }
}
</style>
