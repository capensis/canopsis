<template lang="pug">
  div
    v-layout(row, wrap, align-center)
      v-flex(xs11)
        v-textarea(
          v-field="comment.message",
          v-validate="rules",
          :disabled="!!comment._id",
          :label="$t('modals.createPbehavior.steps.comments.fields.message')",
          :error-messages="errors.collect(messageFieldName)",
          :name="messageFieldName"
        )
      v-flex(xs1)
        v-btn(
          color="error",
          icon,
          @click="$emit('remove')"
        )
          v-icon delete
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'comment',
    event: 'input',
  },
  props: {
    comment: {
      type: Object,
      required: true,
    },
    max: {
      type: Number,
      default: 255,
    },
  },
  computed: {
    rules() {
      return {
        required: true,
        max: this.max,
      };
    },

    messageFieldName() {
      return `${this.comment._id || this.comment.key}_message`;
    },
  },
};
</script>
