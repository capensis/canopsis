<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile.items-per-page-title
        slot(name="title")
          span {{ $t('settings.defaultNumberOfElementsPerPage') }}
    v-container
      v-select.select(
        v-field="value",
        :label="$t('settings.elementsPerPage')",
        :items="$config.PAGINATION_PER_PAGE_VALUES",
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
