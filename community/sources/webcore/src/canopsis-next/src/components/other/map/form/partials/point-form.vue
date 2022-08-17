<template lang="pug">
  v-layout(column)
    c-entity-field(
      v-field="form.entity",
      :label="$t('map.defineEntity')",
      :required="!isLinked",
      :clearable="isLinked"
    )
    c-coordinates-field(v-if="coordinates", v-field="form.coordinates")
    c-enabled-field(v-model="isLinked", :label="$t('map.addLink')")
    c-map-field(v-show="isLinked", v-field="form.map", :required="isLinked")
</template>

<script>
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
