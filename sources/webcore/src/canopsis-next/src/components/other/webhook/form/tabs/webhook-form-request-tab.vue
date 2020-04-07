<template lang="pug">
  div
    v-layout(v-for="(request, index) in requests", :key="index", row)
      v-flex
        webhook-form-request-field(
          v-field="requests[index]",
          :name="`requests[${index}]`",
          :disabled="disabled"
        )
      v-btn(v-if="!disabled", :disabled="isRemoveDisabled", color="error", icon, @click="removeItemFromArray(index)")
        v-icon delete
    v-layout(v-if="!disabled", row)
      v-btn.ma-0.mt-2(color="primary", @click="addRequestHandler") {{ $t('common.add') }}
</template>

<script>
import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formArrayMixin from '@/mixins/form/array';

import TextPairs from '@/components/forms/fields/text-pairs.vue';
import WebhookFormRequestField from '@/components/other/webhook/partials/webhook-form-request-field.vue';

import { getDefaultRequestField } from '@/helpers/forms/webhook';

export default {
  inject: ['$validator'],
  components: { WebhookFormRequestField, TextPairs },
  mixins: [
    formMixin,
    formValidationHeaderMixin,
    formArrayMixin,
  ],
  model: {
    prop: 'requests',
    event: 'input',
  },
  props: {
    requests: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isRemoveDisabled() {
      return this.requests.length === 1;
    },
  },
  methods: {
    addRequestHandler() {
      this.addItemIntoArray(getDefaultRequestField());
    },
  },
};
</script>
