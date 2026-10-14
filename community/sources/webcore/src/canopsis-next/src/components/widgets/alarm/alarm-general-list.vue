<template>
  <c-collapse-panel class="c-alternative-bg-panel" expanded>
    <template #header>
      <span class="font-weight-medium text-uppercase">
        {{ title }}
      </span>
    </template>

    <v-data-table :headers="headers" :items="items" :row-props="rowProps" />
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
    const { t } = useI18n();

    const title = computed(() => {
      const count = props.items.length === 1 ? '' : ` (${props.items.length})`;

      return `${t('common.alarm', props.items.length)}${count}`;
    });

    const headers = computed(() => [
      {
        title: t('common.author'),
        key: 'v.state.a',
        sortable: false,
      },
      {
        title: t('common.connector'),
        key: 'v.connector',
        sortable: false,
      },
      {
        title: t('common.component'),
        key: 'v.component',
        sortable: false,
      },
      {
        title: t('common.resource'),
        key: 'v.resource',
        sortable: false,
      },
    ]);

    const rowProps = ({ item }) => {
      const rowItem = item?.raw ?? item;
      const itemClass = props.itemClass ? rowItem?.[props.itemClass] : undefined;

      return itemClass ? { class: itemClass } : undefined;
    };

    return {
      title,
      headers,
      rowProps,
    };
  },
};
</script>
