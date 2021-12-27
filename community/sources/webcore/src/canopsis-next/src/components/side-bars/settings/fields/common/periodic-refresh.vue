<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div(:class="validationHeaderClass") {{ $t('settings.periodicRefresh') }}
          span.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      periodic-refresh-field(v-field="value", :name="name")
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

export default {
  inject: ['$validator'],
  components: { PeriodicRefreshField },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      required: false,
    },
  },
  created() {
    if (!this.value?.unit) {
      this.updateField('unit', TIME_UNITS.second);
    }
  },
};
</script>
