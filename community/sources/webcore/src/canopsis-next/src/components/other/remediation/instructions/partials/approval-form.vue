<template lang="pug">
  v-layout(column)
    c-description-field(
      v-if="dismissing",
      v-model="form.comment",
      :label="$tc('common.comment')",
      name="dismiss_comment",
      required
    )
    v-layout(justify-end)
      template(v-if="!dismissing")
        v-btn.warning(
          depressed,
          flat,
          @click="showDismissComment"
        ) {{ $t('common.dismiss') }}
        v-btn(
          :loading="submitting",
          :disabled="disabled",
          color="primary",
          @click="approve"
        ) {{ $t('common.approve') }}
      template(v-else)
        v-btn(
          depressed,
          flat,
          @click="cancelDismiss"
        ) {{ $t('common.cancel') }}
        v-btn(
          :loading="submitting",
          :disabled="disabled || errors.any()",
          depressed,
          flat,
          @click="dismiss"
        ) {{ $t('common.dismiss') }}
          v-icon(color="error", right) cancel
</template>

<script>
import { VALIDATION_DELAY } from '@/constants';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  props: {
    submitting: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      form: {
        comment: '',
      },
      dismissing: false,
    };
  },
  watch: {
    dismissing() {
      this.errors.clear();
    },
  },
  methods: {
    closeDismissComment() {
      this.dismissing = false;
    },

    showDismissComment() {
      this.dismissing = true;
    },

    cancelDismiss() {
      this.form.comment = '';
      this.closeDismissComment();
    },

    async dismiss() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.$emit('dismiss', this.form.comment);
      }
    },

    approve() {
      this.$emit('approve');
    },
  },
};
</script>
