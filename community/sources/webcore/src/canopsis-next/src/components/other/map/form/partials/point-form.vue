<template lang="pug">
  v-layout(column)
    c-entity-field(
      v-field="form.entity",
      :label="$t('map.defineEntity')",
      :required="!isLinked",
      :clearable="isLinked",
      :entity-types="entityTypes"
    )
    c-coordinates-field(v-if="coordinates", v-field="form.coordinates")
    c-enabled-field(v-model="isLinked", :label="$t('map.addLink')")
    c-map-field(
      v-show="isLinked",
      v-field="form.map",
      :required="isLinked",
      hide-details
    )
</template>

<script>
import { ENTITY_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    coordinates: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    entityTypes() {
      return Object.values(ENTITY_TYPES);
    },

    isLinked: {
      set(value) {
        this.updateField('map', value ? '' : undefined);
      },
      get() {
        return this.form.map !== undefined && this.form.map !== null;
      },
    },
  },
};
</script>
