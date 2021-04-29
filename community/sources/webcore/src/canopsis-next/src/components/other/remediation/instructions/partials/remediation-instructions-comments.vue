<template lang="pug">
  v-card
    v-list
      v-list-tile(v-if="!comments.length")
        v-list-tile-content
          v-list-tile-title {{ $t('tables.noData') }}
      template(v-for="(comment, index) in comments")
        v-list-tile(:key="index")
          v-list-tile-content
            v-list-tile-title {{ comment.comment }}
          v-list-tile-action
            rating-field(:value="comment.rating", readonly)
        v-divider(v-if="index < comments.length - 1", :key="`divider_${index}`")
</template>

<script>
import { entitiesRemediationInstructionsMixin } from '@/mixins/entities/remediation/instructions';

import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },
  mixins: [entitiesRemediationInstructionsMixin],
  props: {
    remediationInstruction: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      comments: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: comments } = await this.fetchRemediationInstructionCommentsListWithoutStore({
        id: this.remediationInstruction._id,
      });
      this.comments = comments;

      this.pending = false;
    },
  },
};
</script>
