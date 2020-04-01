<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.counterLevels.title') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      v-list-group(data-test="elementsPerPage")
        v-list-tile.items-per-page-title(slot="activator")
          slot(name="title")
            span {{ $t('settings.defaultNumberOfElementsPerPage') }}
        v-container
          v-select.select(
            v-field="form.counter",
            :label="$t('settings.mainCounter')",
            :items="availableCounters",
            hide-details,
            single-line,
            dense
          )
      v-divider
      field-criticity-levels(v-field="form.criticityLevels")
      v-divider
      field-levels-colors-selector(
        v-field="form.criticityLevelsColors",
        colorType="hex",
        hideSuffix
      )
</template>

<script>
import { AVAILABLE_COUNTERS } from '@/constants';

import FieldCriticityLevels from '../fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '../fields/stats/levels-colors-selector.vue';

export default {
  components: {
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableCounters() {
      return Object.values(AVAILABLE_COUNTERS);
    },
  },
};
</script>
