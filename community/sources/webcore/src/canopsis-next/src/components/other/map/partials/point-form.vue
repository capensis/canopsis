<template lang="pug">
  v-layout(column)
    c-entity-field(v-field="form.entity", :label="$t('map.defineEntity')", :required="!isLinked")
    c-coordinate-field(v-if="coordinate", v-field="form.coordinate")
    c-enabled-field(v-model="isLinked", :label="$t('map.addLink')")
    c-map-field(v-show="isLinked", v-field="form.link")
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
    coordinate: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    isLinked: {
      set(value) {
        this.updateField('link', value ? '' : undefined);
      },
      get() {
        return this.form.link !== undefined;
      },
    },
  },
};
</script>
