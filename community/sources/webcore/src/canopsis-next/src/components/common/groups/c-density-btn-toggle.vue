<template lang="pug">
  v-btn-toggle.density-btn-toggle(v-field="value", mandatory)
    v-tooltip(top)
      template(#activator="{ on }")
        v-btn(v-on="on", :value="false", small, flat)
          v-icon(small) $vuetify.icons.density_medium
      span {{ $t('settings.density.comfort') }}
    v-tooltip(top)
      template(#activator="{ on }")
        v-btn(v-on="on", :value="true", small, flat)
          v-icon(small) $vuetify.icons.density_small
      span {{ $t('settings.density.compact') }}
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'change',
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.density-btn-toggle {
  box-shadow: none;
}
</style>
