<template lang="pug">
  v-list-group(data-test="defaultNumberOfElementsPerPage")
    v-list-tile(slot="activator") {{$t('settings.defaultNumberOfElementsPerPage')}}
    v-container(data-test="elementsPerPage")
      v-select.select(
      :label="$t('settings.elementsPerPage')",
      :items="$config.PAGINATION_PER_PAGE_VALUES",
      :value="value",
      @change="$emit('input', $event)",
      hide-details,
      single-line,
      dense,
      )
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

export default {
  props: {
    value: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
  },
};
</script>
