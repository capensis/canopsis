<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div {{ $t('settings.density.title') }}
    v-container
      c-density-btn-toggle(v-field="value", :name="name", column)
</template>

<script>
import { ALARM_DENSE_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: ALARM_DENSE_TYPES.large,
    },
    name: {
      type: String,
      default: 'dense',
    },
  },
};
</script>

<style lang="scss" scoped>
.density {
  &-icon {
    padding: 2px;
    border-radius: 5px;
    border: 1px solid #707070;
  }

  &-radio-group ::v-deep .v-input__control {
    width: 100%;
  }
}
</style>
