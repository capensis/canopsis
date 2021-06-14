<template lang="pug">
  div
    rating-field.justify-center(v-field="form.rating", hover)
    v-textarea(v-field="form.comment", :label="$tc('common.comment')")
</template>

<script>
import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
