<template lang="pug">
  div
    v-layout(row)
      strong {{ $tc('common.comment', comments.length) }}
    v-layout(v-for="(comment, index) in comments", :key="comment.key", row, wrap, allign-center)
      v-flex(xs11)
        v-textarea(
          :disabled="!!comment._id",
          :label="$t('modals.createPbehavior.fields.message')",
          :value="comment.message",
          @input="updateFieldInArrayItem(index, 'message', $event)"
        )
      v-flex(xs1)
        v-btn(color="error", icon, @click="removeItemFromArray(index)")
          v-icon delete
    v-layout(row)
      v-btn.ml-0.primary(type="button", @click="addComment") {{ $t('modals.createPbehavior.buttons.addComment') }}
</template>

<script>
import uid from '@/helpers/uid';

import authMixin from '@/mixins/auth';
import formArrayMixin from '@/mixins/form/array';

export default {
  mixins: [authMixin, formArrayMixin],
  model: {
    prop: 'comments',
    event: 'input',
  },
  props: {
    comments: {
      type: Array,
      required: true,
    },
  },
  methods: {
    addComment() {
      this.addItemIntoArray({
        key: uid(),
        author: this.currentUser._id,
        message: '',
      });
    },
  },
};
</script>
