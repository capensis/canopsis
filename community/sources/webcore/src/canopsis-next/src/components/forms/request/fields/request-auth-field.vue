<template>
  <v-layout class="gap-2">
    <v-flex xs6>
      <c-name-field
        v-field="value.username"
        :label="$t('common.username')"
        :name="`${name}.username`"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex xs6>
      <c-password-field
        v-field="value.password"
        :name="`${name}.password`"
        :disabled="disabled"
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'auth',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
