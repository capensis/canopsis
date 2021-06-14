<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $t('common.description') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-flex(xs12)
          v-card
            v-card-text
              pre {{ pbehaviorException.description }}
    v-tab {{ $t('common.periods') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-flex(xs12)
          v-card
            v-card-text
              pbehavior-exceptions-field(:exdates="preparedPbehaviorExceptionExdates", disabled)
</template>

<script>
import { pbehaviorExceptionToForm } from '@/helpers/forms/exceptions-pbehavior';

import PbehaviorExceptionsField from '../form/pbehavior-exceptions-field.vue';

export default {
  components: { PbehaviorExceptionsField },
  props: {
    pbehaviorException: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    preparedPbehaviorExceptionExdates() {
      const { exdates } = pbehaviorExceptionToForm(this.pbehaviorException);

      return exdates;
    },
  },
};
</script>
