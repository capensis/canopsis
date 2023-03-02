<template lang="pug">
  v-radio-group.density__radio-group(v-if="column", v-field="value", :name="name")
    v-layout.mb-3(v-for="type in types", :key="type.value", row)
      v-flex(xs6)
        v-radio(
          :value="type.value",
          :label="type.text",
          color="primary"
        )
      v-flex(xs6)
        v-icon.density__icon {{ type.icon }}

  v-btn-toggle.density__btn-toggle(v-else, v-field="value", :name="name", mandatory)
    v-tooltip(v-for="type in types", :key="type.value", top)
      template(#activator="{ on }")
        v-btn(v-on="on", :value="type.value", small, flat)
          v-icon(small) {{ type.icon }}
      span {{ type.text }}
</template>

<script>
import { ALARM_DENSE_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'change',
  },
  props: {
    value: {
      type: Number,
      default: ALARM_DENSE_TYPES.large,
    },
    name: {
      type: String,
      default: 'dense',
    },
    column: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    types() {
      return [{
        value: ALARM_DENSE_TYPES.large,
        icon: '$vuetify.icons.density_large',
        text: this.$t('settings.density.comfort'),
      }, {
        value: ALARM_DENSE_TYPES.medium,
        icon: '$vuetify.icons.density_medium',
        text: this.$t('settings.density.compact'),
      }, {
        value: ALARM_DENSE_TYPES.small,
        icon: '$vuetify.icons.density_small',
        text: this.$t('settings.density.ultraCompact'),
      }];
    },
  },
};
</script>

<style lang="scss">
.density {
  &__btn-toggle {
    box-shadow: none;
  }

  &__radio-group .v-input__control {
    width: 100%;
  }

  &__icon {
    padding: 2px;
    border-radius: 5px;
    border: 1px solid #707070;
  }
}
</style>
