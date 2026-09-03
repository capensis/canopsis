<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab :value="TABS.patterns">
        {{ t('common.pattern', 2) }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <v-window-item :value="TABS.patterns">
        <c-layout class="py-3" justify="center">
          <div :style="gridFlexStyle(8)">
            <v-card>
              <v-card-text>
                <tag-patterns-form :form="patterns" readonly />
              </v-card-text>
            </v-card>
          </div>
        </c-layout>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { tagFilterPatternsToForm } from '@/helpers/entities/tag/form';
import { gridFlexStyle } from '@/helpers/layout';

import { useI18n } from '@/hooks/i18n';

import TagPatternsForm from '@/components/other/tag/form/tag-patterns-form.vue';

const TABS = {
  patterns: 'patterns',
};

export default {
  components: {
    TagPatternsForm,
  },

  props: {
    tag: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const activeTab = ref(TABS.patterns);

    const patterns = computed(() => tagFilterPatternsToForm(props.tag));

    return {
      TABS,
      activeTab,
      patterns,
      t,
      gridFlexStyle,
    };
  },
};
</script>
