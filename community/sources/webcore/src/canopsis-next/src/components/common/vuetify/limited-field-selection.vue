<template>
  <template v-if="index < showCount">
    <slot />
  </template>
  <span v-else-if="index === showCount">
    +{{ remaining }} {{ t('common.more') }}
  </span>
  <span v-else class="c-limited-field-selection__hidden" />
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    index: {
      type: Number,
      required: true,
    },
    showCount: {
      type: Number,
      required: true,
    },
    total: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const remaining = computed(() => Math.max(props.total - props.showCount, 0));

    return {
      t,
      remaining,
    };
  },
};
</script>

<style lang="scss">
.v-select__selection:has(.c-limited-field-selection__hidden),
.v-autocomplete__selection:has(.c-limited-field-selection__hidden),
.v-combobox__selection:has(.c-limited-field-selection__hidden) {
  display: none;
}
</style>
