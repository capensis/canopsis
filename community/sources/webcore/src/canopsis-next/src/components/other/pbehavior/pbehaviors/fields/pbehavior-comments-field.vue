<template>
  <div>
    <pbehavior-comment-field
      v-for="(comment, index) in comments"
      v-field="comments[index]"
      :key="comment.key"
      @remove="removeItemFromArray(index)"
    />
    <v-layout>
      <v-btn
        class="ml-0 primary"
        type="button"
        @click="addComment"
      >
        {{ $t('modals.createPbehavior.steps.comments.buttons.addComment') }}
      </v-btn>
    </v-layout>
  </div>
</template>

<script>
import { uid } from '@/helpers/uid';

import { formArrayMixin } from '@/mixins/form';

import PbehaviorCommentField from './pbehavior-comment-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorCommentField },
  mixins: [formArrayMixin],
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
        message: '',
      });
    },
  },
};
</script>
