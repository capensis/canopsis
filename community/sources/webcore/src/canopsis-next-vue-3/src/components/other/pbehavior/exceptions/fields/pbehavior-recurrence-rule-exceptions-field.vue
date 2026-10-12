<template>
  <c-layout gap="3" column>
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
  </c-layout>
</template>

<script>
import PbehaviorExceptionsList from '../../pbehaviors/partials/pbehavior-exceptions-list.vue';

import PbehaviorExceptionsField from './pbehavior-exceptions-field.vue';

export default {
  emits: ['update:exdates'],
  components: {
    PbehaviorExceptionsList,
    PbehaviorExceptionsField,
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
