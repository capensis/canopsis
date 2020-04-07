<template lang="pug">
  div
    v-layout(v-for="(postProcessor, index) in postProcessors", :key="index", row)
      v-flex
        webhook-form-post-processor-field(
          v-field="postProcessors[index]",
          :disabled="disabled",
          :name="`postProcessors[${index}]`"
        )
      v-btn(v-if="!disabled", :disabled="isRemoveDisabled", color="error", icon, @click="removeItemFromArray(index)")
        v-icon delete
    v-layout(row)
      v-btn.ma-0.mt-2(v-if="!disabled", color="primary", @click="addPostProcessorHandler") {{ $t('common.add') }}
</template>

<script>
import TextPairs from '@/components/forms/fields/text-pairs.vue';
import WebhookFormPostProcessorField from '@/components/other/webhook/partials/webhook-form-post-processor-field.vue';

import { getDefaultPostProcessorField } from '@/helpers/forms/webhook';

import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formArrayMixin from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  components: { WebhookFormPostProcessorField, TextPairs },
  mixins: [formArrayMixin, formValidationHeaderMixin],
  model: {
    prop: 'postProcessors',
    event: 'input',
  },
  props: {
    postProcessors: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isRemoveDisabled() {
      return this.postProcessors.length === 1;
    },
  },
  methods: {
    addPostProcessorHandler() {
      this.addItemIntoArray(getDefaultPostProcessorField());
    },
  },
};
</script>
