<template>
  <div>
    <rating-field
      v-field="form.rating"
      class="justify-center"
      hover
    />
    <v-textarea
      v-field="form.comment"
      :label="$tc('common.comment')"
    />
  </div>
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
