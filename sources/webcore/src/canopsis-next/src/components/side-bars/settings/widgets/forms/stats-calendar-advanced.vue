<template lang="pug">
  v-list-group(data-test="advancedSettings")
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ $t('settings.advancedSettings') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-filters(
        :filters="value.filters",
        hideSelect,
        @update:filters="updateField('filters', $event)"
      )
      v-divider
      field-opened-resolved-filter(v-field="value.alarmsStateFilter")
      v-divider
      field-switcher(
        v-field="value.considerPbehaviors",
        :title="$t('settings.considerPbehaviors.title')",
        data-test="widgetConsiderPbehaviors"
      )
      v-divider
      field-criticity-levels(
        v-field="value.criticityLevels",
        data-test="widgetCriticityLevels"
      )
      v-divider
      field-levels-colors-selector(
        v-field="value.criticityLevelsColors",
        colorType="hex",
        hideSuffix
      )
</template>

<script>
import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import FieldOpenedResolvedFilter from '../fields/alarm/opened-resolved-filter.vue';
import FieldFilters from '../fields/common/filters.vue';
import FieldSwitcher from '../fields/common/switcher.vue';
import FieldCriticityLevels from '../fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '../fields/stats/levels-colors-selector.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  inject: ['$validator'],
  components: {
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldSwitcher,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
  },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
