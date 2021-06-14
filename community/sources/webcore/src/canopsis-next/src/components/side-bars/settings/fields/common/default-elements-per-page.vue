<template lang="pug">
  v-list-group(data-test="elementsPerPage")
    v-list-tile.items-per-page-title(slot="activator")
      slot(name="title")
        span {{ $t('settings.defaultNumberOfElementsPerPage') }}
    v-container(data-test="elementsPerPageFieldContainer")
      v-select.select(
        :label="$t('settings.elementsPerPage')",
        :items="$config.PAGINATION_PER_PAGE_VALUES",
        :value="value",
        @change="$emit('input', $event)",
        hide-details,
        single-line,
        dense
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

<style lang="scss" scoped>
  .items-per-page-title {
    & /deep/ .v-list__tile.theme--light {
      padding-right: 4px !important;
    }
  }
</style>
