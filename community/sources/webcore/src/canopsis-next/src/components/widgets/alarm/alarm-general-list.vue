<template>
  <c-collapse-panel class="c-alternative-bg-panel" expanded>
    <template #header>
      <span class="font-weight-medium text-uppercase">{{ title }}</span>
    </template>
    <v-data-table
      :headers="headers"
      :items="items"
      :item-class="itemClass"
    />
  </c-collapse-panel>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    itemClass: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();

    const title = computed(() => {
      const count = props.items.length === 1 ? '' : ` (${props.items.length})`;

      return `${tc('common.alarm', props.items.length)}${count}`;
    });

    const headers = computed(() => [
      {
        text: t('common.author'),
        value: 'v.state.a',
        sortable: false,
      },
      {
        text: t('common.connector'),
        value: 'v.connector',
        sortable: false,
      },
      {
        text: t('common.component'),
        value: 'v.component',
        sortable: false,
      },
      {
        text: t('common.resource'),
        value: 'v.resource',
        sortable: false,
      },
    ]);

    return {
      title,
      headers,
    };
  },
};
</script>
