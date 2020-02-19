<template lang="pug">
  v-list-group(data-test="periodicRefresh")
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{$t('settings.periodicRefresh')}}
        span.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      periodic-refresh-field(v-field="value")
</template>

<script>
import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

export default {
  inject: ['$validator'],
  components: { PeriodicRefreshField },
  mixins: [formMixin, formValidationHeaderMixin],
  props: {
    value: {
      type: Object,
      default: () => ({ ...DEFAULT_PERIODIC_REFRESH }),
    },
  },
  created() {
    if (this.value && !this.value.unit) {
      this.updateField('unit', DEFAULT_PERIODIC_REFRESH.unit);
    }
  },
};
</script>
