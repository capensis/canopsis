<template lang="pug">
  v-form(@submit.prevent.stop="submit")
    v-card(width="400")
      v-card-title.primary.pa-2.white--text
        v-layout(justify-space-between, align-center)
          h4 {{ title || $t('mermaid.addPoint') }}
          v-btn.ma-0.ml-3(
            icon,
            small,
            @click="close"
          )
            v-icon(color="white") close
      point-form(v-field="form")
      v-layout(justify-end)
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="close"
        ) {{ $t('common.cancel') }}
        v-btn.error(
          v-if="removable",
          :disabled="submitting",
          depressed,
          flat,
          @click="$emit('remove')"
        ) {{ $t('common.delete') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    submittableMixinCreator(),
    confirmableModalMixinCreator({
      closeMethod: 'close',
    }),
  ],
  props: {
    point: {
      type: Object,
      required: true,
    },
    title: {
      type: String,
      required: false,
    },
    removable: {
      type: Boolean,
      required: false,
    },
  },
  data() {
    return {
      form: { ...this.point },
    };
  },
  methods: {
    close() {
      this.$emit('cancel');
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.$emit('submit', this.form);
      }
    },
  },
};
</script>
