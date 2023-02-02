<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div {{ $t('settings.density.title') }}
    v-container
      v-radio-group.density-radio-group(v-field="value", :name="name")
        v-layout.mb-3(row)
          v-flex(xs6)
            v-radio(
              :value="false",
              :label="$t('settings.density.comfort')",
              color="primary"
            )
          v-flex(xs6)
            v-icon.density-icon $vuetify.icons.density_medium
        v-layout(row)
          v-flex(xs6)
            v-radio(
              :value="true",
              :label="$t('settings.density.compact')",
              color="primary"
            )
          v-flex(xs6)
            v-icon.density-icon $vuetify.icons.density_small
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
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
