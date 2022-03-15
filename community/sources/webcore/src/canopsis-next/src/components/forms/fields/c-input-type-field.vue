<template lang="pug">
  v-select.c-input-type-field(
    v-field="value",
    :items="preparedTypes",
    :label="label",
    :disabled="disabled",
    :flat="flat",
    :error-messages="errorMessages",
    hide-details,
    dense
  )
    template(slot="selection", slot-scope="{ parent, item, index }")
      v-icon.c-input-type-field__icon(small) {{ getInputTypeIcon(item.value) }}
    template(slot="item", slot-scope="{ item }")
      v-list-tile-avatar.c-input-type-field__avatar
        v-icon.c-input-type-field__icon(small) {{ getInputTypeIcon(item.value) }}
      v-list-tile-content
        v-list-tile-title {{ item.text }}
</template>

<script>
import { PATTERN_INPUT_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    flat: {
      type: Boolean,
      default: false,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    types: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    preparedTypes() {
      return this.types.map(
        type => (type.text ? type : ({ ...type, text: this.$t(`mixedField.types.${type.value}`) })),
      );
    },
  },
  methods: {
    getInputTypeIcon(type) {
      return {
        [PATTERN_INPUT_TYPES.string]: 'title',
        [PATTERN_INPUT_TYPES.number]: 'exposure_plus_1',
        [PATTERN_INPUT_TYPES.boolean]: 'toggle_on',
        [PATTERN_INPUT_TYPES.null]: 'space_bar',
        [PATTERN_INPUT_TYPES.array]: 'view_array',
      }[type];
    },
  },
};
</script>

<style lang="scss">
.c-input-type-field {
  &__icon {
    color: inherit !important;
    opacity: .6;
  }

  &__avatar {
    min-width: 30px;

    & .v-avatar {
      width: 20px !important;
      height: 20px !important;
    }
  }
}
</style>
