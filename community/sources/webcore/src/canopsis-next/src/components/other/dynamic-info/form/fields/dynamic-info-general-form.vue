<template>
  <div>
    <c-id-field
      v-field="form._id"
      :disabled="isDisabledIdField"
    />
    <c-name-field
      v-field="form.name"
      required
    />
    <c-enabled-field v-field="form.enabled" />
    <c-disable-during-periods-field v-field="form.disable_during_periods" />
    <c-description-field
      v-field="form.description"
      required
    />
  </div>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
