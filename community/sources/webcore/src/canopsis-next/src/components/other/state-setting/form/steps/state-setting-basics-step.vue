<template lang="pug">
  v-layout(column)
    v-layout.mb-2(row)
      v-flex(xs7)
        c-name-field.mr-2(
          v-field="form.title",
          :label="$t('common.title')",
          name="title",
          required
        )
      v-flex(xs3)
        c-priority-field.mx-2(
          v-field="form.priority",
          required
        )
      v-flex(xs2)
        c-enabled-field.ml-2(v-field="form.enabled")
    state-setting-method-field(v-field="form.method")
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

import StateSettingMethodField from '../fields/state-setting-method-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingMethodField },
  mixins: [formValidationHeaderMixin],
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
};
</script>

<style lang="scss">
.state-setting-form {
  background-color: transparent !important;
}
</style>
