<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.defaultNumberOfElementsPerPage')}}
    v-container
      v-select.select(
      :label="$t('settings.elementsPerPage')",
      :items="items",
      v-model="value",
      @change="$emit('input', $event)"
      hide-details,
      single-line,
      dense,
      )
</template>

<script>
import { PAGINATION_PER_PAGE_VALUES } from '@/config';

export default {
  props: {
    value: {
      type: Number,
      default: PAGINATION_PER_PAGE_VALUES[0],
    },
  },
  data() {
    return {
      items: PAGINATION_PER_PAGE_VALUES,
    };
  },
};
</script>
