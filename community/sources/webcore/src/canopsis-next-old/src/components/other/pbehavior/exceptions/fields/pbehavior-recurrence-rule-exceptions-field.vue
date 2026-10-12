<template>
  <v-layout class="gap-3" column>
    <pbehavior-exceptions-list
      v-if="exceptions.length"
      :exceptions="exceptions"
      :disabled="disabled"
    />
    <pbehavior-exceptions-field
      :exdates="exdates"
      :disabled="disabled"
      :with-exdate-type="withExdateType"
    />
  </v-layout>
</template>

<script>
import PbehaviorExceptionsList from '../../pbehaviors/partials/pbehavior-exceptions-list.vue';

import PbehaviorExceptionsField from './pbehavior-exceptions-field.vue';

export default {
  components: {
    PbehaviorExceptionsList,
    PbehaviorExceptionsField,
  },
  model: {
    prop: 'exdates',
    event: 'input',
  },
  props: {
    exdates: {
      type: Array,
      default: () => [],
    },
    exceptions: {
      type: Array,
      default: () => [],
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
