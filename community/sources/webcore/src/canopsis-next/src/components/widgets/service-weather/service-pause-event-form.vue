<template>
  <div>
    <c-description-field
      v-field="form.comment"
      :label="$tc('common.comment')"
      name="comment"
      required
    />
    <c-pbehavior-reason-field
      v-field="form.reason"
      required
      return-object
    />
  </div>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
