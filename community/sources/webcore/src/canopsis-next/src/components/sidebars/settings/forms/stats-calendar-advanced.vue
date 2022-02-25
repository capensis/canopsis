<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ $t('settings.advancedSettings') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-filters(
        :filters="value.filters",
        hide-select,
        addable,
        editable,
        @update:filters="updateField('filters', $event)"
      )
      v-divider
      field-opened-resolved-filter(v-field="value.opened")
      v-divider
      field-switcher(
        v-field="value.considerPbehaviors",
        :title="$t('settings.considerPbehaviors.title')"
      )
      v-divider
      field-criticity-levels(v-field="value.criticityLevels")
      v-divider
      field-levels-colors-selector(
        v-field="value.criticityLevelsColors",
        color-type="hex",
        hide-suffix
      )
</template>

<script>
import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import FieldOpenedResolvedFilter from '@/components/sidebars/settings/fields/alarm/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import FieldCriticityLevels from '@/components/sidebars/settings/fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/settings/fields/stats/levels-colors-selector.vue';

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
