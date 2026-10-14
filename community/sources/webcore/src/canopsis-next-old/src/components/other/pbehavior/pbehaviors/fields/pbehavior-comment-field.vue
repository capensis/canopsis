<template>
  <div>
    <v-layout align-center>
      <v-textarea
        v-field="comment.message"
        v-validate="rules"
        :disabled="!!comment._id"
        :label="$t('common.message')"
        :error-messages="errors.collect(messageFieldName)"
        :name="messageFieldName"
      />
      <c-action-btn
        type="delete"
        @click="$emit('remove')"
      />
    </v-layout>
  </div>
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
