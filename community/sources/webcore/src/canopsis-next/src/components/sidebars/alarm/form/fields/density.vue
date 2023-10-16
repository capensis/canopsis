<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>
        <div>{{ $t('settings.density.title') }}</div>
      </v-list-item>
    </template>
    <v-container>
      <c-density-btn-toggle
        v-field="value"
        :name="name"
        column
      />
    </v-container>
  </v-list-group>
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
