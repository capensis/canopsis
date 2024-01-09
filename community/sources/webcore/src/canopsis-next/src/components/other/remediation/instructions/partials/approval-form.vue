<template lang="pug">
  v-layout(column)
    c-description-field(
      v-if="dismissing",
      v-model="form.comment",
      :label="$tc('common.comment')",
      name="comment"
    )
    v-layout(v-if="!disabled", justify-end)
      template(v-if="!dismissing")
        v-btn.warning(
          depressed,
          flat,
          @click="showDismissComment"
        ) {{ $t('common.dismiss') }}
        v-btn(
          :loading="submitting",
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
          depressed,
          flat,
          @click="dismiss"
        ) {{ $t('common.dismiss') }}
          v-icon(color="error", right) cancel
</template>

<script>
export default {
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

    dismiss() {
      this.$emit('dismiss', this.form.comment);
    },

    approve() {
      this.$emit('approve');
    },
  },
};
</script>
