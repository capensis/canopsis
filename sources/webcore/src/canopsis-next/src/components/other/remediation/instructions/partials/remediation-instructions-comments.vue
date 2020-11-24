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
import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },
  props: {
    comments: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
