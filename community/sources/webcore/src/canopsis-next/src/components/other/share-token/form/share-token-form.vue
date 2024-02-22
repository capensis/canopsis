<template>
  <v-layout column>
    <c-name-field
      v-field="form.description"
      required
    />
    <c-information-block :title="$t('shareToken.tokenExpiration')">
      <c-enabled-field v-field="form.duration_enabled" />
      <c-duration-field
        v-if="form.duration_enabled"
        v-field="form.duration"
        long
      />
    </c-information-block>
  </v-layout>
</template>

<script>
export default {
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
